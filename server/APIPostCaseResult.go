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
	req, err := t.entities.DB.DB().Prepare(`UPDATE Cases SET in_progress=false, passed=false, finished=true WHERE id=?`)
	check(err)
	_, err = req.Exec(caseID)
	check(err)

	//run hook
	sessionID, cmdLine := t.getSessionIDandCMDLineByCaseID(caseID)
	t.hook_FailInSession(sessionID, cmdLine)
}

func (t *CombatServer) markCasePassed(caseID string) {
	req, err := t.entities.DB.DB().Prepare(`UPDATE Cases SET in_progress=false, passed=true, finished=true WHERE id=?`)
	check(err)
	_, err = req.Exec(caseID)
	check(err)
}

func (t *CombatServer) markCaseNotInProgress(caseID string) {
	req, err := t.entities.DB.DB().Prepare(`UPDATE Cases SET in_progress=false WHERE id=?`)
	check(err)
	_, err = req.Exec(caseID)
	check(err)
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

func (t *CombatServer) hook_first_failInSession(sessionID string, cmdLine string) {
	fmt.Println("Hook_FirstFailInSession" + sessionID + " " + cmdLine)
	go t.alarmSlack_FirstFailInSession(sessionID, cmdLine)
}

// The method set first fail flag as true, if it was false, and return true for first call.
func (t *CombatServer) IsFirstFailInSession(sessionID string) bool {
	req, err := t.entities.DB.DB().Prepare(`UPDATE Sessions SET hook_first_fail=true WHERE id=? AND hook_first_fail=false`) // Set FirstFail flag as true, if not true yet
	check(err)
	execRes, err := req.Exec(sessionID)
	check(err)
	rowsAffected, err := execRes.RowsAffected()
	check(err)
	if rowsAffected != 0 { // if it first fail in the session
		return true
	} else {
		return false
	}
}

func (t *CombatServer) hook_FailInSession(sessionID string, cmdLine string) {
	fmt.Println("Hook_FailInSession: " + sessionID + " " + cmdLine)
	if t.IsFirstFailInSession(sessionID) {
		t.hook_first_failInSession(sessionID, cmdLine)
	}
}

func (t *CombatServer) getSessionIDandCMDLineByCaseID(caseID string) (string, string) {
	req, err := t.entities.DB.DB().Prepare(`SELECT session_id,cmd_line FROM Cases WHERE id=?`) // get session_idcmd_linene for alarm
	check(err)
	rows, err := req.Query(caseID)
	check(err)
	rows.Next()
	var sessionID string
	var cmdLine string
	check(rows.Scan(&sessionID, &cmdLine))
	check(rows.Close())
	return sessionID, cmdLine
}

func (t *CombatServer) IsTryOutFalseNegative(stdOut string) bool {
	for _, curPattern := range t.config.FalseNegativePatterns {
		r, err := regexp.Compile(curPattern)
		check(err)
		if r.MatchString(stdOut) {
			return true
		}
	}
	return false
}

func (t *CombatServer) setCaseResultHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	check(err)

	reqStruct, err := apireqresp.ParseReqPostCaseResultFromBytes(body)
	check(err)

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

	req, err := t.entities.DB.DB().Prepare(`SELECT id FROM tries WHERE case_id=?`) // get count of tries
	check(err)
	rows, err := req.Query(caseID)
	check(err)
	triesCount := 0
	for rows.Next() {
		triesCount++
	}
	check(rows.Close())

	fmt.Println("CurrentTryCount=" + strconv.Itoa(triesCount))

	req, err = t.entities.DB.DB().Prepare("INSERT INTO tries(case_id,exit_status,std_out) VALUES(?,?,?)")
	check(err)
	res, err := req.Exec(caseID, exitStatus, stdOut)
	check(err)
	tryID64, err := res.LastInsertId()
	check(err)
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

	if exitStatus != "0" {
		check(os.MkdirAll("./tries/"+tryID, 0777))

		decodedFile, err := reqStruct.GetDecodedFile()
		check(err)

		f, err := os.OpenFile("./tries/"+tryID+"/out_archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
		check(err)

		_, err = io.Copy(f, bytes.NewReader(decodedFile))
		check(err)
		check(f.Close())

		go unzip("./tries/"+tryID+"/out_archived.zip", "./tries/"+tryID)
	}

	fmt.Println(r.RemoteAddr + " Provide result for case: " + caseID + ". Status=" + exitStatus)

}
