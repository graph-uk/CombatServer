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
		//fmt.Println(r.RemoteAddr + " Getting a job")
		//select case not in progress
		var caseID, caseCMD, sessionID string
		rows, err := t.entities.DB.DB().Query(`SELECT id, cmd_line, session_id FROM Cases WHERE finished=false AND in_progress=false ORDER BY RANDOM() LIMIT 1`)
		check(err)

		// if found some case not in progress.
		if rows.Next() {
			check(rows.Scan(&caseID, &caseCMD, &sessionID))
			check(rows.Close())

			// set case.InProgress = true, and unlock DB
			curTime := time.Now()
			req, err := t.entities.DB.DB().Prepare("UPDATE Cases SET in_progress=?, started_at=? WHERE id=?")
			check(err)

			_, err = req.Exec(true, curTime, caseID)
			check(err)

			zipFile, err := ioutil.ReadFile("./sessions/" + sessionID + "/archived.zip")
			check(err)

			resp := apireqresp.NewResGetJob(caseID, caseCMD, zipFile)

			respJson, err := resp.GetJson()
			check(err)

			_, err = w.Write(respJson)
			check(err)

			fmt.Println(r.RemoteAddr + " Get a job (CasesRun) for case: " + caseCMD)
		} else { // when not found cases to run
			rows.Close()
			//fmt.Println(r.RemoteAddr + " Jobs not found. Idle. ")
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
