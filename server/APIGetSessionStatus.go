package server

import (
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"time"
	"strconv"
)

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

		if totalCasesCount == finishedCasesCount && totalCasesCount != 0 {
			w.Header().Set("Finished", "True")
			if failedCasesCount == 0 {
				w.Write([]byte("Success. All tests passed"))
			} else {
				w.Write([]byte("Finished. Errors: " + strconv.Itoa(failedCasesCount) + "\r\n"))
				for _, curCase := range errorCases {
					w.Write([]byte("    " + curCase + "\r\n"))
				}
			}
		} else {
			w.Header().Set("Finished", "False")
			if totalCasesCount != 0 {
				if failedCasesCount != 0 {
					w.Write([]byte("Running (" + strconv.Itoa(finishedCasesCount) + "/" + strconv.Itoa(totalCasesCount) + ") Errors: " + strconv.Itoa(failedCasesCount) + "\r\n"))
					for _, curCase := range errorCases {
						w.Write([]byte("    " + curCase + "\r\n"))
					}
				} else {
					w.Write([]byte("Running (" + strconv.Itoa(finishedCasesCount) + "/" + strconv.Itoa(totalCasesCount) + ")"))
				}
			} else {
				w.Write([]byte("Cases exploring"))
			}
		}

		fmt.Println(r.Host + " Get session status: for session: " + sessionID)
	}
}
