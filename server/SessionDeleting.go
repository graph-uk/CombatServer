package server

import (
	"fmt"
	"os"

	"github.com/graph-uk/combat-server/server/config"
)

func (t *CombatServer) DeleteOldSessions() {
	oldSessions := t.GetOldSessionsList()
	for _, curSessionID := range oldSessions {
		t.DeleteSession(*curSessionID)
		//fmt.Println(*curSessionID)
	}
}

func (t *CombatServer) GetOldSessionsList() []*string {

	req, err := t.mdb.DB.Prepare(`SELECT id FROM Sessions ORDER BY id DESC`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	rows, err := req.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var result []*string
	for rows.Next() {
		var curSessionID string
		rows.Scan(&curSessionID)
		result = append(result, &curSessionID)
	}
	rows.Close()

	if config.GetApplicationConfig().MaxStoredSessions < 1 { // when count less than 1 - sessions stores forever (no old sessions)
		result = []*string{}
	} else {
		if len(result) >= config.GetApplicationConfig().MaxStoredSessions {
			result = result[config.GetApplicationConfig().MaxStoredSessions:]
		} else {
			result = []*string{}
		}
	}
	return result
}

func (t *CombatServer) DeleteSession(sessionID string) {
	slash := string(os.PathSeparator)
	fmt.Println("Deleting session started. SessionID=" + sessionID)

	req, err := t.mdb.DB.Prepare(`SELECT id FROM Cases WHERE sessionID=?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rows, err := req.Query(sessionID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var cases []*string // Cases collecting
	for rows.Next() {
		var curCase string
		rows.Scan(&curCase)
		cases = append(cases, &curCase)
	}
	rows.Close()

	for _, curCaseID := range cases {
		fmt.Println("Deleting case: " + *curCaseID)
		req, err := t.mdb.DB.Prepare(`SELECT id FROM Tries WHERE caseID=?`)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		rows, err := req.Query(*curCaseID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var tries []*string // Tries collecting
		for rows.Next() {
			var curTry string
			rows.Scan(&curTry)
			tries = append(tries, &curTry)
		}
		rows.Close()

		for _, curTryID := range tries {
			fmt.Println("Deleting try: " + *curTryID)
			os.RemoveAll(t.startPath + slash + "tries" + slash + *curTryID) // delete tries from directory
		}

		req, err = t.mdb.DB.Prepare(`DELETE FROM Tries WHERE caseID=?`) // delete tries from DB
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		req.Exec(*curCaseID)

		os.RemoveAll(t.startPath + slash + "cases" + slash + sessionID) // delete session from directory

		req, err = t.mdb.DB.Prepare(`DELETE FROM Cases WHERE sessionID=?`) // delete cases from DB
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		req.Exec(sessionID)

	}

	req, err = t.mdb.DB.Prepare(`DELETE FROM Sessions WHERE id=?`) // delete session from DB
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Exec(sessionID)

}
