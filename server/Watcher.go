package server

import (
	"fmt"
	"time"

	"github.com/graph-uk/combat-server/data/models/status"

	"github.com/graph-uk/combat-server/data/repositories"

	"github.com/graph-uk/combat-server/utils"
)

func checkNotificationEnabled(config *utils.Config) {
	configsRepo := &repositories.Configs{}
	dbConfig := configsRepo.Find()
	if !dbConfig.NotificationEnabled {
		fmt.Println(time.Now().Sub(dbConfig.MuteTimestamp))
		if time.Now().After(dbConfig.MuteTimestamp.Add(time.Duration(config.NotificationMuteDurationMinutes) * time.Minute)) { // if mute time has left
			dbConfig.NotificationEnabled = true
			configsRepo.Update(dbConfig)
		}
	}
}

func checkCases() {
	casesRepo := &repositories.Cases{}

	cases := casesRepo.FindProcessingCases()
	currentTime := time.Now()

	for _, sessionCase := range cases {
		if sessionCase.DateStarted.Add(time.Duration(utils.GetApplicationConfig().CaseTimeoutSec) * time.Second).Before(currentTime) {
			sessionCase.Status = status.Pending
			casesRepo.Update(&sessionCase)
			fmt.Printf("Watcher: case %d is out of date. Restarted.\n", sessionCase.ID)
		}
	}
}

// TimeoutWatcher ...
func TimeoutWatcher(config *utils.Config) {
	for {
		checkCases()
		checkNotificationEnabled(config)
		time.Sleep(10 * time.Second)
	}
}
