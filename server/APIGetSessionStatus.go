package server

import (
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"time"
	//"strconv"
	"encoding/json"
)

type SessionStatus struct {
	Finished           bool
	TotalCasesCount    int
	FinishedCasesCount int
	FailReports        []string
}

func (t *CombatServer) getSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//sessionID := r.Header.Get("sessionID")
		r.ParseMultipartForm(32 << 20)
		sessionID := r.FormValue("sessionID")
		if sessionID == "" {
			fmt.Println("cannot extract session ID")
			return
		}

		req, err := t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=?`)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err := req.Query(sessionID)
		if err != nil {
			fmt.Println(err)
			return
		}
		var totalCasesCount int
		rows.Next()
		err = rows.Scan(&totalCasesCount)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows.Close()

		req, err = t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=? AND finished="true"`)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err = req.Query(sessionID)
		if err != nil {
			fmt.Println(err)
			return
		}
		var finishedCasesCount int
		rows.Next()
		err = rows.Scan(&finishedCasesCount)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows.Close()

		req, err = t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=? AND finished="true" AND passed="false"`)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err = req.Query(sessionID)
		if err != nil {
			fmt.Println(err)
			return
		}
		var failedCasesCount int
		rows.Next()
		err = rows.Scan(&failedCasesCount)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows.Close()

		req, err = t.mdb.DB.Prepare(`SELECT cmdLine FROM Cases WHERE sessionID=? AND finished="true" AND passed="false"`)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err = req.Query(sessionID)
		if err != nil {
			fmt.Println(err)
			return
		}
		var errorCases []string
		for rows.Next() {
			var cmdLine string
			err = rows.Scan(&cmdLine)
			if err != nil {
				fmt.Println(err)
				return
			}
			errorCases = append(errorCases, cmdLine)
		}
		rows.Close()

		var sessionStatus SessionStatus
		sessionStatus.TotalCasesCount = totalCasesCount
		sessionStatus.FinishedCasesCount = finishedCasesCount
		if totalCasesCount == finishedCasesCount && totalCasesCount != 0 {
			sessionStatus.Finished = true
		} else {
			sessionStatus.Finished = false
		}

		for _, curCase := range errorCases {
			sessionStatus.FailReports = append(sessionStatus.FailReports, curCase)
		}
		sessionStatusJSON, _ := json.Marshal(sessionStatus)
		w.Write(sessionStatusJSON)
	}
}
