package server

import (
	"fmt"
	"strconv"
	"time"
)

//func markCaseNotInProgress(caseID string) {

//	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress="false" WHERE id=?`)
//	check(err)
//	_, err = req.Exec(caseID)
//	check(err)
//}

func (t *CombatServer) CheckCases(curtime time.Time) {

	rows, err := t.mdb.DB.Query(`SELECT id,startedAt FROM Cases WHERE (startedAt IS NOT NULL) AND (inProgress="1") AND (finished="false")`)
	if err != nil {
		fmt.Println(err)
		return
	}
	var oldRunCases []int
	for rows.Next() {
		//fmt.Println("sdf")
		var startedAt time.Time
		var id int
		rows.Scan(&id, &startedAt)
		if startedAt.Add(100 * time.Second).Before(curtime) {
			oldRunCases = append(oldRunCases, id)
			//			req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress="false" WHERE id=?`)
			//			check(err)
			//			_, err = req.Exec(id)
			//			check(err)
			//			fmt.Println("Watcher: case " + id + " is out of date. Restart.")
		}
	}
	rows.Close()
	for _, curID := range oldRunCases {
		t.markCaseNotInProgress(strconv.Itoa(curID))
		fmt.Println("Watcher: case " + strconv.Itoa(curID) + " is out of date. Restarted.")
	}
}

func (t *CombatServer) TimeoutWatcher() {
	for {
		curTime := time.Now()
		t.mdb.Lock()
		t.CheckCases(curTime)
		fmt.Println("Watcher: Cases checked")
		t.mdb.Unlock()
		time.Sleep(10 * time.Second)
	}
}
