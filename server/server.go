package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	//	"github.com/graph-uk/combat-server/server/DB"
	"github.com/graph-uk/combat-server/server/config"
	"github.com/graph-uk/combat-server/server/entities"
)

type CombatServer struct {
	config    *config.Config
	startPath string
	//db        DB.DB
	entities *entities.Entities
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkFolder(folderName string) {
	if _, err := os.Stat(folderName); os.IsNotExist(err) { // if folder does not exist - try to create
		check(os.MkdirAll(folderName, 0777))
	}
}

func checkFolders() {
	checkFolder("sessions")
	checkFolder("tries")
}

func NewCombatServer() *CombatServer {
	var result CombatServer
	var err error
	result.startPath, err = os.Getwd()
	check(err)
	result.config, err = config.LoadConfig()
	check(err)
	//	check(result.db.Connect("./base.sl3?_busy_timeout=60000"))
	//	result.db.CheckDBNew()
	result.entities = entities.NewEntities(`./base.sl3?_busy_timeout=60000`)
	checkFolders()
	return &result
}

func (t *CombatServer) Serve() {
	go t.TimeoutWatcher()
	http.Handle("/tries/", http.StripPrefix("/tries/", http.FileServer(http.Dir("./tries"))))

	http.HandleFunc("/api/v1/sessions", t.createSessionHandler)
	http.HandleFunc("/api/v1/sessions/", t.getSessionStatusHandler)
	http.HandleFunc("/api/v1/commands/get-job", t.getJobHandler)
	http.HandleFunc("/api/v1/commands/setCaseResult", t.setCaseResultHandler)
	http.HandleFunc("/sessions/", t.pageSessionStatusHandler)

	fmt.Println("Serving combat tests at port: " + strconv.Itoa(t.config.Port) + "...")
	check(http.ListenAndServe(":"+strconv.Itoa(t.config.Port), nil))
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
