package server

import (
	"fmt"
	"time"

	"github.com/graph-uk/combat-server/data/models/status"

	"github.com/graph-uk/combat-server/data/repositories"

	"github.com/graph-uk/combat-server/server/config"
)

func checkCases() {
	casesRepo := &repositories.Cases{}

	cases := casesRepo.FindProcessingCases()
	currentTime := time.Now()

	for _, sessionCase := range cases {
		if sessionCase.DateStarted.Add(time.Duration(config.GetApplicationConfig().CaseTimeoutSec) * time.Second).Before(currentTime) {
			sessionCase.Status = status.Pending
			casesRepo.Update(&sessionCase)
			fmt.Printf("Watcher: case %d is out of date. Restarted.\n", sessionCase.ID)
		}
	}
}

// TimeoutWatcher ...
func TimeoutWatcher() {
	for {
		checkCases()
		time.Sleep(10 * time.Second)
	}
}
