package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/graph-uk/combat-server/server/apireqresp"
)

func (t *CombatServer) markCaseFailed(caseID string) {
	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress=false, passed=false, finished=true WHERE id=?`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = req.Exec(caseID)
	if err != nil {
		fmt.Println(err)
		return
	}

	//run hook
	sessionID, cmdLine, err := t.getSessionIDandCMDLineByCaseID(caseID)
	if err != nil {
		fmt.Println(err)
		t.mdb.Unlock()
		return
	}
	t.hook_FailInSession(sessionID, cmdLine)
}

func (t *CombatServer) markCasePassed(caseID string) {
	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress=false, passed=true, finished=true WHERE id=?`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = req.Exec(caseID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (t *CombatServer) markCaseNotInProgress(caseID string) {
	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress=false WHERE id=?`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = req.Exec(caseID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (t *CombatServer) alarmSlack_FirstFailInSession(sessionID string, cmdLine string) {
	if t.config.FirstSessionFailSlackHook == "" {
		fmt.Println("Slack hook is not set. Slack message is not sent.")
		return
	}

	if t.config.FirstSessionFailSlackChannel == "" {
		fmt.Println("Slack channel is not set. Slack message is not sent.")
		return
	}

	alarmMessage := `{
		"channel": "#` + t.config.FirstSessionFailSlackChannel + `",
	    "text": "<` + t.config.ServerHostname + ":" + strconv.Itoa(t.config.Port) + `/sessions/` + sessionID + `|` + t.config.ProjectName + ` testing failed>: _` + cmdLine + `_"
		}`

	bodyBuffer := new(bytes.Buffer)
	bodyBuffer.WriteString(alarmMessage)

	fmt.Println(alarmMessage)

	resp, err := http.Post("https://hooks.slack.com/services/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", "application/json", bodyBuffer)
	if err != nil {
		fmt.Println("Cannot post message. Message may be was not sent.")
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read response body. Message may be was not sent.")
		return
	}

	if string(respBody) == "ok" {
		fmt.Println("Slack message sent.")
	} else {
		fmt.Println("Slack message sending failed with response:")
		fmt.Println(string(respBody))
	}
}

func (t *CombatServer) hook_FirstFailInSession(sessionID string, cmdLine string) {
	fmt.Println("Hook_FirstFailInSession" + sessionID + " " + cmdLine)
	go t.alarmSlack_FirstFailInSession(sessionID, cmdLine)
}

// The method set first fail flag as true, if it was false, and return true for first call.
func (t *CombatServer) IsFirstFailInSession(sessionID string) (bool, error) {
	req, err := t.mdb.DB.Prepare(`UPDATE Sessions SET hook_FirstFail=true WHERE id=? AND hook_FirstFail=false`) // Set FirstFail flag as true, if not true yet
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	execRes, err := req.Exec(sessionID)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	rowsAffected, err := execRes.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if rowsAffected != 0 { // if it first fail in the session
		return true, nil
	} else {
		return false, nil
	}
}

func (t *CombatServer) hook_FailInSession(sessionID string, cmdLine string) {
	fmt.Println("Hook_FailInSession: " + sessionID + " " + cmdLine)

	ItFirstFail, err := t.IsFirstFailInSession(sessionID)
	if err != nil {
		return
	}

	if ItFirstFail {
		t.hook_FirstFailInSession(sessionID, cmdLine)
	}
}

func (t *CombatServer) getSessionIDandCMDLineByCaseID(caseID string) (string, string, error) {
	req, err := t.mdb.DB.Prepare(`SELECT sessionID,cmdLine FROM Cases WHERE id=?`) // get sessionID, cmdLine for alarm
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	rows, err := req.Query(caseID)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	rows.Next()
	var sessionID string
	var cmdLine string
	err = rows.Scan(&sessionID, &cmdLine)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	rows.Close()
	return sessionID, cmdLine, nil
}

func (t *CombatServer) IsTryOutFalseNegative(stdOut string) bool {
	for _, curPattern := range t.config.FalseNegativePatterns {
		r, _ := regexp.Compile(curPattern)
		if r.MatchString(stdOut) {
			return true
		}
	}
	return false
}

func (t *CombatServer) setCaseResultHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	reqStruct, err := apireqresp.ParseReqPostCaseResultFromBytes(body)
	if err != nil {
		fmt.Println(err)
	}

	caseID := reqStruct.CaseID
	if caseID == "" {
		fmt.Println("cannot extract caseID")
		return
	}

	exitStatus := reqStruct.ExitStatus
	if exitStatus == "" {
		fmt.Println("cannot extract exitStatus")
		return
	}

	stdOut := reqStruct.StdOut
	if stdOut == "" {
		fmt.Println("cannot extract stdOut")
		return
	}

	if t.IsTryOutFalseNegative(stdOut) { // drop false-negative result.
		t.markCaseNotInProgress(caseID)
		fmt.Println("False-negative dropped")
		return
	}

	t.mdb.Lock()
	req, err := t.mdb.DB.Prepare(`SELECT id FROM Tries WHERE caseID=?`) // get count of tries
	if err != nil {
		fmt.Println(err)
		t.mdb.Unlock()
		return
	}
	rows, err := req.Query(caseID)
	if err != nil {
		fmt.Println(err)
		t.mdb.Unlock()
		return
	}
	triesCount := 0
	for rows.Next() {
		triesCount++
	}
	rows.Close()

	fmt.Println("CurrentTryCount=" + strconv.Itoa(triesCount))

	req, err = t.mdb.DB.Prepare("INSERT INTO Tries(caseID,exitStatus,stdOut) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
		t.mdb.Unlock()
		return
	}
	res, err := req.Exec(caseID, exitStatus, stdOut)
	if err != nil {
		fmt.Println(err)
		t.mdb.Unlock()
		return
	}
	tryID64, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		t.mdb.Unlock()
		return
	}
	tryID := strconv.Itoa(int(tryID64))

	fmt.Println("TestTriesCount=" + strconv.Itoa(triesCount))
	if exitStatus == "0" { // if test passed
		t.markCasePassed(caseID)
	} else { // if test failed
		if triesCount+2 > t.config.CountOfRetries { // if test failed too many times
			t.markCaseFailed(caseID)
		} else { // if test failed, and should try again
			t.markCaseNotInProgress(caseID)
		}
	}

	t.mdb.Unlock()

	if exitStatus != "0" {
		os.MkdirAll("./tries/"+tryID, 0777)

		decodedFile, err := reqStruct.GetDecodedFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := os.OpenFile("./tries/"+tryID+"/out_archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		io.Copy(f, bytes.NewReader(decodedFile))
		f.Close()

		go unzip("./tries/"+tryID+"/out_archived.zip", "./tries/"+tryID)
	}

	fmt.Println(r.RemoteAddr + " Provide result for case: " + caseID + ". Status=" + exitStatus)

}
