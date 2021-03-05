package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"malibu-server/data/models/status"

	"malibu-server/data/repositories"

	"malibu-server/utils"
)

func check(err error) {
	if err != nil {
		//panic(err)
		log.Println(err)
	}
}

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

func deleteSessionsItemsOutOfRange(config *utils.Config) {
	sessionsRepo := &repositories.Sessions{}
	sessionsRepo.DeleteOldSessions(config.MaxStoredSessions)
}

// func deleteSuccessArtifacts() {
// 	triesRepo := &repositories.Tries{}
// 	allTries := triesRepo.FindAll()
// 	for _, curTry := range allTries {
// 		if curTry.ExitStatus == `0` {
// 			os.RemoveAll(fmt.Sprintf(`_data/tries/%d`, curTry.ID))
// 		}
// 	}
// }

// func deleteOldSuccessfullRuns() {
// 	files, err := ioutil.ReadDir(`_data/tries-succ`)

// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	for _, file := range files {
// 		if strings.HasPrefix(file.Name(), `old`) {
// 			os.RemoveAll(`_data/tries-succ/` + file.Name())
// 		}
// 	}
// }

// func deleteSessionsOutOfRange(config *utils.Config) {
// 	sessionsRepo := &repositories.Sessions{}
// 	sessionsRepo.DeleteOldSessions(config.MaxStoredSessions)
// }

func removeSessionsFoldersOlderThanLast(sessionsFoldersList []string, lastSessionID string) {
	lastSessionIDInt, err := strconv.Atoi(lastSessionID)
	if err != nil { //
		log.Println(`removeSessionsOlderThanLast: lastSessionID cannot be casted to int: "` + lastSessionID + `"`)
		return
	}
	for _, curSessionFolder := range sessionsFoldersList {
		curSessionIDInt, err := strconv.Atoi(curSessionFolder)
		if err != nil {
			continue
		}
		if curSessionIDInt < lastSessionIDInt {
			//log.Println(`removeSessionsFoldersOlderThanLast ` + fmt.Sprintf(`_data/sessions/%d`, curSessionFolder))
			os.RemoveAll(fmt.Sprintf(`_data/sessions/%d`, curSessionFolder))
		}
	}
}

func removeZippedTries(files []string) {
	for _, file := range files {
		//log.Println(`removeZippedTries ` + fmt.Sprintf(`_data/tries/%d/archived.zip`, file))
		os.RemoveAll(fmt.Sprintf(`_data/tries/%d/archived.zip`, file))
	}
}

func getLastSession() string {
	sessionsRepo := &repositories.Sessions{}
	lastSession := sessionsRepo.FindLast()
	if lastSession != nil {
		return lastSession.ID
	}
	return ""
}

func getTriesFoldersList() []string {
	files, err := ioutil.ReadDir(`_data/tries`)
	check(err)
	res := []string{}
	for _, file := range files {
		if file.IsDir() {
			res = append(res, file.Name())
		}
	}
	return res
}

func getSessionsFoldersList() []string {
	files, err := ioutil.ReadDir(`_data/sessions`)
	check(err)
	res := []string{}
	for _, file := range files {
		if file.IsDir() {
			res = append(res, file.Name())
		}
	}
	return res
}

func getOldSuccessfulRuns() []string {
	files, err := ioutil.ReadDir(`_data/tries-succ`)
	check(err)
	res := []string{}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), `old`) {
			res = append(res, file.Name())
		}
	}
	return res
}

func deleteTriesFoldersAndItemsWithNoCaseItem(triesFlodersList []string) {
	casesRepo := &repositories.Cases{}
	triesRepo := &repositories.Tries{}
	//	allCases := casesRepo.FindAll()
	for _, curTryID := range triesFlodersList {
		tryIDInt, err := strconv.Atoi(curTryID)
		if err != nil {
			continue
		}
		try := triesRepo.Find(tryIDInt)
		if try == nil { // the try's item is not found by ID
			os.RemoveAll(fmt.Sprintf(`_data/tries/%d`, curTryID))
		}
		caseItem := casesRepo.Find(try.CaseID)
		if caseItem == nil {
			triesRepo.DeleteByID(tryIDInt)
		}
	}
}

func deleteCasesItemsWithNoSessionItem() {
	casesRepo := &repositories.Cases{}
	sessionsRepo := &repositories.Sessions{}
	allCases := casesRepo.FindAll()
	for _, curCase := range allCases {
		session := sessionsRepo.Find(curCase.SessionID)
		if session == nil {
			casesRepo.DeleteByID(curCase.ID)
		}
	}
}

func clearTmpFolder() {
	files, err := ioutil.ReadDir(`/tmp`)
	//	check(err)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), `combatOutput`) {
			os.RemoveAll(fmt.Sprintf(`/tmp`, file.Name()))
		}
	}
}

func removeOldData(config *utils.Config) {
	sessionsFlodersList := getSessionsFoldersList()
	triesFlodersList := getTriesFoldersList()
	lastSession := getLastSession()
	oldSuccessfullRuns := getOldSuccessfulRuns()

	if true {
		log.Println("sessionsFolder: ", sessionsFlodersList)
		log.Println("triesFolder: ", triesFlodersList)
		log.Println("lastSession: ", lastSession)
		log.Println("oldSuccessfulRuns: ", oldSuccessfullRuns)
	}
	removeZippedTries(triesFlodersList)
	removeSessionsFoldersOlderThanLast(sessionsFlodersList, lastSession)
	deleteSessionsItemsOutOfRange(config)
	deleteCasesItemsWithNoSessionItem()
	deleteTriesFoldersAndItemsWithNoCaseItem(triesFlodersList)
	clearTmpFolder()
}

// TimeoutWatcher ...
func TimeoutWatcher(config *utils.Config) {
	divider := 10
	curIteration := 0
	for {
		checkCases()
		checkNotificationEnabled(config)
		//deleteSuccessArtifacts()
		// deleteOldSuccessfullRuns()
		// deleteSessionsOutOfRange(config)
		//		if curIteration == divider {
		if curIteration < 0 {
			curIteration = divider
			removeOldData(config)
		}
		time.Sleep(10 * time.Second)
		curIteration--
	}
}
