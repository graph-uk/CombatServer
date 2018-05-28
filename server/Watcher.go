package server

import (
	"fmt"
	"strconv"
	"time"
)

func (t *CombatServer) CheckCases() {
	curtime := time.Now()
	rows, err := t.mdb.DB.Query(`SELECT id,startedAt FROM Cases WHERE (startedAt IS NOT NULL) AND (inProgress=true) AND (finished=false)`)
	if err != nil {
		fmt.Println(err)
		return
	}
	var oldRunCases []int
	for rows.Next() {
		var startedAt time.Time
		var id int
		rows.Scan(&id, &startedAt)
		if startedAt.Add(time.Duration(t.config.CaseTimeoutSec) * time.Second).Before(curtime) {
			oldRunCases = append(oldRunCases, id)
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
		t.mdb.Lock()
		t.CheckCases()
		t.mdb.Unlock()
		time.Sleep(10 * time.Second)
	}
}
