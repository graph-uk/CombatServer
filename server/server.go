package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/graph-uk/combat-server/server/config"
	"github.com/graph-uk/combat-server/server/mutexedDB"
)

type CombatServer struct {
	config    *config.Config
	startPath string
	mdb       mutexedDB.MutexedDB
}

func checkFolder(folderName string) error {
	if _, err := os.Stat(folderName); os.IsNotExist(err) { // if folder does not exist - try to create
		err := os.MkdirAll(folderName, 0777)
		if err != nil {
			fmt.Println("Cannot create folder " + folderName)
			return err
		}
	} else {
		err := os.MkdirAll(folderName+string(os.PathSeparator)+"TMP_TESTFOLDER", 0777)
		if err != nil {
			fmt.Println("Cannot create subfolder in folder " + folderName + ". Check permissions")
			return err
		}

		err = os.RemoveAll(folderName + string(os.PathSeparator) + "TMP_TESTFOLDER")
		if err != nil {
			fmt.Println("Cannot delete subfolder in folder " + folderName + ". Check permissions")
			return err
		}
	}
	return nil
}

func checkFolders() error {
	err := checkFolder("sessions")
	if err != nil {
		return err
	}
	err = checkFolder("tries")
	if err != nil {
		return err
	}

	return nil
}

func NewCombatServer() (*CombatServer, error) {
	var result CombatServer
	var err error
	result.startPath, err = os.Getwd()
	result.config, err = config.LoadConfig()
	if err != nil {
		return &result, err
	}

	err = result.mdb.Connect("./base.sl3")
	if err != nil {
		return &result, err
	}

	err = checkFolders()
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func (t *CombatServer) Serve() error {
	go t.TimeoutWatcher()

	//http.Handle("/tries/", http.FileServer(http.Dir(t.startPath+string(os.PathSeparator)+"tries")))

	//http.Handle("/tries/", http.FileServer(http.Dir("./tries")))
	//http.Handle("/tries/", http.StripPrefix("/tries/", http.FileServer(http.Dir("./tries"))))
	http.Handle("/tries/", http.StripPrefix("/tries/", http.FileServer(http.Dir("./tries"))))

	http.HandleFunc("/createSession", t.createSessionHandler)
	http.HandleFunc("/getJob", t.getJobHandler)
	http.HandleFunc("/setCaseResult", t.setCaseResultHandler)
	http.HandleFunc("/getSessionStatus", t.getSessionStatusHandler)
	http.HandleFunc("/sessions/", t.pageSessionStatusHandler)

	err := http.ListenAndServe(":"+strconv.Itoa(t.config.Port), nil)
	return err
}

func (t *CombatServer) addToGOPath(pathExtention string) []string {
	result := os.Environ()
	for curVarIndex, curVarValue := range result {
		if strings.HasPrefix(curVarValue, "GOPATH=") {
			result[curVarIndex] = result[curVarIndex] + string(os.PathListSeparator) + pathExtention
			return result
		}
	}
	return result
}
