package server

import (
	"fmt"
	//"io"
	"io/ioutil"
	"net/http"
	"time"
)

func (t *CombatServer) getJobHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t.mdb.Lock()
		//defer t.mdb.Unlock()

		var caseID, caseCMD, sessionID string
		rows, err := t.mdb.DB.Query(`SELECT id, cmdLine, sessionID FROM Cases WHERE finished="false" AND inProgress="false" ORDER BY RANDOM() LIMIT 1`)
		if err != nil {
			fmt.Println(err)
			return
		}

		if rows.Next() {
			err = rows.Scan(&caseID, &caseCMD, &sessionID)
			if err != nil {
				fmt.Println(err)
				t.mdb.Unlock()
				return
			}
			rows.Close()

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

			w.Header().Add("Command", "RunCase")
			w.Header().Add("Params", caseCMD)
			w.Header().Add("SessionID", caseID)

			zipArchive, err := ioutil.ReadFile("./sessions/" + sessionID + "/archived.zip")
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = w.Write(zipArchive)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(r.RemoteAddr + " Get a job (CasesRun) for case: " + caseCMD)

		} else { // when not found cases to run
			rows.Close()
			w.Header().Add("Command", "idle")
			t.mdb.Unlock()
		}
	}
}
