package server

import (
	"net/http"
	//"strconv"
	"encoding/json"
	"strings"

	"github.com/graph-uk/combat-server/server/session"
)

//type SessionStatus struct {
//	Finished           bool
//	//TotalCasesCount    int
//	//FinishedCasesCount int
//	//FailReports        []string
//}

func (t *CombatServer) pageSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		path := strings.Split(r.URL.Path, "/")
		sessionID := strings.TrimSpace(path[len(path)-1])
		if strings.TrimSpace(sessionID) == "" {
			w.Write([]byte("Session ID is not specified. Please, provide session ID like: /sessions/11203487203498"))
			return
		}
		w.Write([]byte("Session: " + sessionID + "<br>\n"))

		// Get session status
		session, err := session.NewAssignedSession(sessionID, &t.mdb)
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		casesArray, err := session.GetCasesArray()
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		json, err := json.Marshal(casesArray)
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		w.Write(json)

		//		sessionStatus, err := session.GetStatus()
		//		if err != nil {
		//			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
		//			return
		//		}

		//		// Print session status as HTML
		//		finishedStr := "False"
		//		if sessionStatus.Finished {
		//			finishedStr = "True"
		//		}

		//		w.Write([]byte("Finished: " + finishedStr + "<br>\n"))
		//		w.Write([]byte("Params: " + sessionStatus.Params + "<br>\n"))
		//		w.Write([]byte("TotalCases: " + strconv.Itoa(sessionStatus.TotalCasesCount) + "<br>\n"))
		//		w.Write([]byte("FinishedCases: " + strconv.Itoa(sessionStatus.FinishedCasesCount) + "<br>\n"))
	}
}
