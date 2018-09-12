package server

import (
	"fmt"
	"os"
)

func (t *CombatServer) DeleteOldSessions() {
	oldSessions := t.GetOldSessionsList()
	for _, curSessionID := range oldSessions {
		t.DeleteSession(*curSessionID)
	}
}

func (t *CombatServer) GetOldSessionsList() []*string {

	req, err := t.entities.DB.DB().Prepare(`SELECT id FROM Sessions ORDER BY id DESC`)
	check(err)
	rows, err := req.Query()
	check(err)

	var result []*string
	for rows.Next() {
		var curSessionID string
		check(rows.Scan(&curSessionID))
		result = append(result, &curSessionID)
	}
	check(rows.Close())

	if t.config.CountOfStoredSessions < 1 { // when count less than 1 - sessions stores forever (no old sessions)
		result = []*string{}
	} else {
		if len(result) >= t.config.CountOfStoredSessions {
			result = result[t.config.CountOfStoredSessions:]
		} else {
			result = []*string{}
		}
	}
	return result
}

func (t *CombatServer) DeleteSession(sessionID string) {
	slash := string(os.PathSeparator)
	fmt.Println("Deleting session started. SessionID=" + sessionID)

	req, err := t.entities.DB.DB().Prepare(`SELECT id FROM Cases WHERE session_id=?`)
	check(err)
	rows, err := req.Query(sessionID)
	check(err)

	var cases []*string // Cases collecting
	for rows.Next() {
		var curCase string
		check(rows.Scan(&curCase))
		cases = append(cases, &curCase)
	}
	check(rows.Close())

	for _, curCaseID := range cases {
		fmt.Println("Deleting case: " + *curCaseID)
		req, err := t.entities.DB.DB().Prepare(`SELECT id FROM tries WHERE case_id=?`)
		check(err)
		rows, err := req.Query(*curCaseID)
		check(err)

		var tries []*string // Tries collecting
		for rows.Next() {
			var curTry string
			check(rows.Scan(&curTry))
			tries = append(tries, &curTry)
		}
		check(rows.Close())

		for _, curTryID := range tries {
			fmt.Println("Deleting try: " + *curTryID)
			check(os.RemoveAll(t.startPath + slash + "tries" + slash + *curTryID)) // delete tries from directory
		}

		req, err = t.entities.DB.DB().Prepare(`DELETE FROM tries WHERE case_id=?`) // delete tries from DB
		check(err)
		_, err = req.Exec(*curCaseID)
		check(err)

		check(os.RemoveAll(t.startPath + slash + "cases" + slash + sessionID)) // delete session from directory

		req, err = t.entities.DB.DB().Prepare(`DELETE FROM Cases WHERE session_id=?`) // delete cases from DB
		check(err)
		_, err = req.Exec(sessionID)
		check(err)
	}

	req, err = t.entities.DB.DB().Prepare(`DELETE FROM Sessions WHERE id=?`) // delete session from DB
	check(err)
	_, err = req.Exec(sessionID)
	check(err)

}
