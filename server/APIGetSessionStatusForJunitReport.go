package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type jUnitCase struct {
	id       int
	cmdLine  string
	finished bool
	Passed   bool
	STDOut   string
}

type jUnitTest struct {
	Name  string
	Cases []*jUnitCase
}

type SessionStatusForJunitReport struct {
	Finished bool
	Tests    []*jUnitTest
}

func (t *CombatServer) getSessionStatusForJunitReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		sessionID := r.FormValue("sessionID")
		if sessionID == "" {
			fmt.Println("cannot extract session ID")
			return
		}

		//getting current session cases
		req, err := t.mdb.DB.Prepare(`SELECT id,cmdLine,finished,passed FROM Cases WHERE sessionID=?`)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err := req.Query(sessionID)
		if err != nil {
			fmt.Println(err)
			return
		}

		var cases []*jUnitCase

		for rows.Next() {
			var Case jUnitCase
			rows.Scan(&Case.id, &Case.cmdLine, &Case.finished, &Case.Passed)
			cases = append(cases, &Case)
		}
		rows.Close()

		//get stdOut for each case
		for curCaseIdx, curCase := range cases {
			req, err := t.mdb.DB.Prepare(`SELECT stdOut FROM tries WHERE caseID=? ORDER BY id DESC LIMIT 1`)
			if err != nil {
				fmt.Println(err)
				return
			}
			rows, err := req.Query(curCase.id)
			if err != nil {
				fmt.Println(err)
				return
			}

			rows.Next()
			rows.Scan(&cases[curCaseIdx].STDOut)
			rows.Close()
			//fmt.Println(curCase.id)
			//fmt.Println(cases[curCaseIdx].STDOut)
		}

		//get session status is finished
		sessionFinished := true
		for _, curCase := range cases {
			if !curCase.finished {
				sessionFinished = false
			}
		}
		//fmt.Println(sessionFinished)

		//get all distinct test names
		jUntTests := map[string]bool{}
		for _, curCase := range cases {
			testname := strings.Split(curCase.cmdLine, ` `)[0]
			if !jUntTests[testname] {
				jUntTests[testname] = true
			}
		}

		//build the result struct
		var SessionStatusForJunitReport SessionStatusForJunitReport
		SessionStatusForJunitReport.Finished = sessionFinished

		for testname, _ := range jUntTests {
			//fmt.Println(testname)
			var test jUnitTest
			test.Name = testname
			for _, curCase := range cases {
				caseTestName := strings.Split(curCase.cmdLine, ` `)[0]
				if caseTestName == testname {
					test.Cases = append(test.Cases, curCase)
				}
			}
			SessionStatusForJunitReport.Tests = append(SessionStatusForJunitReport.Tests, &test)
		}

		//send the result struct
		SessionStatusForJunitReportJSON, _ := json.Marshal(SessionStatusForJunitReport)
		w.Write(SessionStatusForJunitReportJSON)
	}
}
