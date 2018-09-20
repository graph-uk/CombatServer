package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/graph-uk/combat-server/server/config"
)

func (t *CombatServer) CheckCases() {
	curtime := time.Now()
	rows, err := t.mdb.DB.Query(`SELECT id,date_started FROM Cases WHERE (date_started IS NOT NULL) AND (status=1)`)
	if err != nil {
		fmt.Println(err)
		return
	}
	var oldRunCases []int
	for rows.Next() {
		var date_started time.Time
		var id int
		rows.Scan(&id, &date_started)
		if date_started.Add(time.Duration(config.GetApplicationConfig().CaseTimeoutSec) * time.Second).Before(curtime) {
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
