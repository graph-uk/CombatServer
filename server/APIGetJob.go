package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/graph-uk/combat-server/server/apireqresp"
)

func (t *CombatServer) getJobHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t.mdb.Lock()
		//defer t.mdb.Unlock()

		//select case not in progress
		var caseID, caseCMD, sessionID string
		rows, err := t.mdb.DB.Query(`SELECT id, cmdLine, sessionID FROM Cases WHERE finished=false AND inProgress=false ORDER BY RANDOM() LIMIT 1`)
		if err != nil {
			fmt.Println(err)
			t.mdb.Unlock()
			return
		}

		// if found some case not in progress.
		if rows.Next() {
			err = rows.Scan(&caseID, &caseCMD, &sessionID)
			if err != nil {
				fmt.Println(err)
				t.mdb.Unlock()
				return
			}
			rows.Close()

			// set case.InProgress = true, and unlock DB
			curTime := time.Now()
			req, err := t.mdb.DB.Prepare("UPDATE Cases SET inProgress=?, startedAt=? WHERE id=?")
			if err != nil {
				fmt.Println(err)
				t.mdb.Unlock()
				return
			}
			_, err = req.Exec(true, curTime, caseID)
			if err != nil {
				fmt.Println(err)
				t.mdb.Unlock()
				return
			}
			t.mdb.Unlock()

			zipFile, err := ioutil.ReadFile("./sessions/" + sessionID + "/archived.zip")
			if err != nil {
				fmt.Println(err)
				return
			}

			resp := apireqresp.NewResGetJob(caseID, caseCMD, zipFile)

			respJson, err := resp.GetJson()
			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = w.Write(respJson)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(r.RemoteAddr + " Get a job (CasesRun) for case: " + caseCMD)
		} else { // when not found cases to run
			w.WriteHeader(http.StatusNotFound)
			rows.Close()
			t.mdb.Unlock()
		}
	}
}
