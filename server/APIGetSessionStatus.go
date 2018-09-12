package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/graph-uk/combat-server/server/apireqresp"
)

func (t *CombatServer) getSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		itemsPathArr := strings.Split(r.URL.Path, `/`)
		if len(itemsPathArr) < 1 {
			fmt.Println("cannot extract session ID, because items arr has zero length")
			return
		}

		sessionID := itemsPathArr[len(itemsPathArr)-1]

		if sessionID == "" {
			fmt.Println("cannot extract session ID, because it is empty")
			return
		}

		var casesExploringMessage string
		var sessionStatus apireqresp.ResGetSession

		req, err := t.entities.DB.DB().Prepare(`SELECT cases_exploring_fail_message FROM Sessions WHERE id=?`)
		check(err)
		rows, err := req.Query(sessionID)
		check(err)

		fmt.Println(rows.Next())

		rows.Scan(&casesExploringMessage)
		check(rows.Close())

		req, err = t.entities.DB.DB().Prepare(`SELECT Count()as count FROM Cases WHERE session_id=?`)
		check(err)

		rows, err = req.Query(sessionID)
		check(err)

		var totalCasesCount int
		rows.Next()
		if rows.Scan(&totalCasesCount) == nil {
			rows.Close()

			var finishedCasesCount int
			req, err = t.entities.DB.DB().Prepare(`SELECT Count()as count FROM Cases WHERE session_id=? AND finished=true`)
			check(err)
			rows, err = req.Query(sessionID)
			check(err)
			rows.Next()
			check(rows.Scan(&finishedCasesCount))
			check(rows.Close())

			var failedCasesCount int
			req, err = t.entities.DB.DB().Prepare(`SELECT Count()as count FROM Cases WHERE session_id=? AND finished=true AND passed=false`)
			check(err)
			rows, err = req.Query(sessionID)
			check(err)
			rows.Next()
			check(rows.Scan(&failedCasesCount))
			check(rows.Close())

			var errorCases []string
			req, err = t.entities.DB.DB().Prepare(`SELECT cmd_line FROM Cases WHERE session_id=? AND finished=true AND passed=false`)
			check(err)
			rows, err = req.Query(sessionID)
			check(err)
			for rows.Next() {
				var cmdLine string
				check(rows.Scan(&cmdLine))
				errorCases = append(errorCases, cmdLine)
			}
			check(rows.Close())

			sessionStatus.CasesExploringFailMessage = casesExploringMessage
			sessionStatus.TotalCasesCount = totalCasesCount
			sessionStatus.FinishedCasesCount = finishedCasesCount
			if totalCasesCount == finishedCasesCount && totalCasesCount != 0 {
				sessionStatus.Finished = true
			} else {
				sessionStatus.Finished = false
			}
			if sessionStatus.CasesExploringFailMessage != "" {
				sessionStatus.Finished = true
			}

			for _, curCase := range errorCases {
				sessionStatus.FailReports = append(sessionStatus.FailReports, curCase)
			}
		} else {
			rows.Close()
		}

		sessionStatusJSON, err := sessionStatus.GetJson()
		check(err)
		w.Write(sessionStatusJSON)
	}
}
