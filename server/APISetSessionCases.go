package server

import (
	"fmt"
	//"io"
	"net/http"
	//"os"
	"strconv"
	"strings"

	//"time"
)

func (t *CombatServer) setSessionCasesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		r.ParseMultipartForm(32 << 20)

		sessionID := r.FormValue("sessionID")
		if sessionID == "" {
			fmt.Println("cannot extract session ID")
			return
		}

		sessionCases := r.FormValue("cases")
		if sessionCases == "" {
			fmt.Println("cannot extract session cases")
			return
		}

		sessionCasesArr := strings.Split(sessionCases, "\n")

		req, err := t.mdb.DB.Prepare("INSERT INTO Cases(cmdline, sessionID) VALUES(?,?)")
		if err != nil {
			fmt.Println(err)
			return
		}

		casesCount := 0
		for _, curCase := range sessionCasesArr {
			curCaseCleared := strings.TrimSpace(curCase)
			if curCaseCleared != "" {
				casesCount++
				_, err = req.Exec(curCase, sessionID)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		fmt.Println(r.Host + " Provided " + strconv.Itoa(casesCount) + " cases for session: " + sessionID)

	}
}
