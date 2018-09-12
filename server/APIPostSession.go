package server

import (
	"bytes"
	//	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/graph-uk/combat-server/server/apireqresp"
)

func (t *CombatServer) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionName := strconv.FormatInt(time.Now().UnixNano(), 10)

	body, err := ioutil.ReadAll(r.Body)
	check(err)
	reqStruct, err := apireqresp.ParseReqPostSessionFromBytes(body)
	check(err)

	//create dir, and save the file.
	check(os.MkdirAll("./sessions/"+sessionName, 0777))
	f, err := os.OpenFile("./sessions/"+sessionName+"/archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
	check(err)
	defer f.Close()
	decodedFile, err := reqStruct.GetDecodedFile()
	check(err)
	_, err = io.Copy(f, bytes.NewReader(decodedFile))
	check(err)

	//create session in DB.
	req, err := t.entities.DB.DB().Prepare("INSERT INTO Sessions(id,params) VALUES(?,?)")
	check(err)
	_, err = req.Exec(sessionName, reqStruct.SessionParams)
	check(err)

	//Mark all unfinished cases as finished and failed
	req, err = t.entities.DB.DB().Prepare(`UPDATE Cases SET in_progress=false, passed=false, finished=true WHERE finished=false`)
	check(err)
	_, err = req.Exec()
	check(err)

	w.Header().Add(`Location`, sessionName)
	w.WriteHeader(http.StatusCreated)

	fmt.Println(r.RemoteAddr + " Create new session: " + sessionName + " " + reqStruct.SessionParams)

	go t.doCasesExplore(reqStruct.SessionParams, sessionName)
}

func (t *CombatServer) doCasesExplore(params, sessionID string) {
	err := unzip("./sessions/"+sessionID+"/archived.zip", "./sessions/"+sessionID+"/unarch")
	check(err)
	check(os.Chdir("./sessions/" + sessionID + "/unarch/src/Tests"))
	rootTestsPath, err := os.Getwd()
	check(err)
	rootTestsPath += string(os.PathSeparator) + ".." + string(os.PathSeparator) + ".."

	command := []string{"cases"}
	for _, curParameter := range strings.Split(params, " ") {
		if strings.TrimSpace(curParameter) != "" {
			command = append(command, curParameter)
		}
	}
	cmd := exec.Command("combat", command...)
	cmd.Env = t.addToGOPath(rootTestsPath)

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	exitStatus := cmd.Run()

	check(os.Chdir(t.startPath))
	if exitStatus == nil {
		t.setCasesForSession(out.String(), sessionID)
	} else {

		req, err := t.entities.DB.DB().Prepare("UPDATE Sessions SET cases_exploring_fail_message=? WHERE id=?")
		check(err)
		_, err = req.Exec(out.String()+"\r\n"+outErr.String(), sessionID) // cases ExploringStarted
		check(err)

		fmt.Println("Cannot extract cases")
		fmt.Println(out.String())
		fmt.Println(outErr.String())
	}
}

func (t *CombatServer) setCasesForSession(sessionCases, sessionID string) {
	sessionCasesArr := strings.Split(sessionCases, "\n")

	req, err := t.entities.DB.DB().Prepare("INSERT INTO Cases(cmd_line, session_id) VALUES(?,?)")
	check(err)

	casesCount := 0
	for _, curCase := range sessionCasesArr {
		curCaseCleared := strings.TrimSpace(curCase)
		if curCaseCleared != "" {
			casesCount++
			_, err = req.Exec(curCase, sessionID)
			check(err)
		}
	}

	fmt.Println("Explored " + strconv.Itoa(casesCount) + " cases for session: " + sessionID)

	go t.DeleteOldSessions()
}
