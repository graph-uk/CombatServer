package server

import (
	"fmt"
	//"io"
	"io/ioutil"
	"net/http"
	"time"
)

func sendJobToNode(w http.ResponseWriter, r *http.Request) {
}

func (t *CombatServer) getJobHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t.mdb.Lock()
		rows, err := t.mdb.DB.Query("SELECT id,params FROM Sessions WHERE status=0 limit 1")
		if err != nil {
			fmt.Println(err)
			return
		}

		if rows.Next() {
			var sessionId, sessionParams string
			err = rows.Scan(&sessionId, &sessionParams)
			if err != nil {
				fmt.Println(err)
				return
			}
			rows.Close()

			req, err := t.mdb.DB.Prepare("UPDATE Sessions SET status=? WHERE id=?")
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = req.Exec(1, sessionId) // cases ExploringStarted
			if err != nil {
				fmt.Println(err)
				return
			}
			t.mdb.Unlock()
			w.Header().Add("Command", "CasesExplore")
			w.Header().Add("Params", sessionParams)
			w.Header().Add("SessionID", sessionId)

			zipArchive, err := ioutil.ReadFile("./sessions/" + sessionId + "/archived.zip")
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = w.Write(zipArchive)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(r.Host + " Get a job (CasesExplore) for session: " + sessionId)
		} else { // when no one session needed to explore cases, check are we have cases to run
			rows.Close()
			var caseID, caseCMD, sessionID string
			rows, err := t.mdb.DB.Query(`SELECT id, cmdLine, sessionID FROM cases WHERE finished="false" AND inProgress="false" ORDER BY RANDOM() LIMIT 1`)
			if err != nil {
				fmt.Println(err)
				return
			}
			if rows.Next() {
				err = rows.Scan(&caseID, &caseCMD, &sessionID)
				if err != nil {
					fmt.Println(err)
					return
				}
				rows.Close()

				curTime := time.Now()
				req, err := t.mdb.DB.Prepare("UPDATE Cases SET inProgress=?, startedAt=? WHERE id=?")
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = req.Exec(true, curTime, caseID)
				if err != nil {
					fmt.Println(err)
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
				fmt.Println(r.Host + " Get a job (CasesRun) for case: " + caseCMD)

			} else { // when not found cases to run
				rows.Close()
				w.Header().Add("Command", "idle")
				t.mdb.Unlock()
			}
		}
	}
}
