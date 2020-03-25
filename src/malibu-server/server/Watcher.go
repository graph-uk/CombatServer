package server

import (
	"fmt"
	"time"

	"malibu-server/data/models/status"

	"malibu-server/data/repositories"

	"malibu-server/utils"
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

func deleteSessionsOutOfRange(config *utils.Config) {
	sessionsRepo := &repositories.Sessions{}
	sessionsRepo.DeleteOldSessions(config.MaxStoredSessions)

}

// TimeoutWatcher ...
func TimeoutWatcher(config *utils.Config) {
	for {
		checkCases()
		checkNotificationEnabled(config)
		deleteSessionsOutOfRange(config)
		time.Sleep(10 * time.Second)
	}
}
