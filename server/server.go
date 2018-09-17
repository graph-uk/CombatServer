package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/graph-uk/combat-server/server/config"
	"github.com/graph-uk/combat-server/server/mutexedDB"
	"github.com/graph-uk/combat-server/server/site"
	"github.com/graph-uk/combat-server/server/site/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	err = result.mdb.Connect("./base.sl3?_busy_timeout=60000")
	if err != nil {
		return &result, err
	}

	err = checkFolders()
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func parseTemplates() (*template.Template, error) {
	cleanRoot := filepath.Clean("server/site/")
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := strings.Replace(path[pfx:], "\\", "/", -1)

			t := root.New(name)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

// Start web server
func (t *CombatServer) Start() error {
	go t.TimeoutWatcher()

	templates, _ := parseTemplates()

	renderer := &site.Template{
		Templates: templates}

	e := echo.New()
	e.Renderer = renderer
	e.Static("/assets", "./assets/_")
	e.Use(middleware.Logger())

	e.GET("/sessions/", sessions.Index)
	e.GET("/sessions/:id", sessions.View)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(t.config.Port)))

	// http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("server/assets"))))
	// http.Handle("/tries/", http.StripPrefix("/tries/", http.FileServer(http.Dir("./tries"))))

	// http.HandleFunc("/getJob", t.getJobHandler)
	// http.HandleFunc("/setCaseResult", t.setCaseResultHandler)
	// http.HandleFunc("/getSessionStatus", t.getSessionStatusHandler)
	// http.HandleFunc("/getSessionStatusForJunitReport", t.getSessionStatusForJunitReportHandler)

	// http.HandleFunc("/createSession", t.createSessionHandler)
	// http.HandleFunc("/sessions/", sessions.Handler)

	// fmt.Println("Serving combat tests at port: " + strconv.Itoa(t.config.Port) + "...")
	// err := http.ListenAndServe(":"+strconv.Itoa(t.config.Port), nil)
	return nil
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
