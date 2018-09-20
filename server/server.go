package server

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/graph-uk/combat-server/server/api/jobs"
	sessionsAPI "github.com/graph-uk/combat-server/server/api/sessions"
	"github.com/graph-uk/combat-server/server/api/tries"
	"github.com/graph-uk/combat-server/server/config"
	"github.com/graph-uk/combat-server/server/mutexedDB"
	"github.com/graph-uk/combat-server/server/site"
	"github.com/graph-uk/combat-server/server/site/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// CombatServer ...
type CombatServer struct {
	startPath string
	mdb       mutexedDB.MutexedDB
}

// NewCombatServer ...
func NewCombatServer() (*CombatServer, error) {
	var result CombatServer
	var err error
	result.startPath, err = os.Getwd()
	if err != nil {
		return &result, err
	}

	err = result.mdb.Connect("./base.sl3?_busy_timeout=60000")
	if err != nil {
		return &result, err
	}

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
	go TimeoutWatcher()

	templates, _ := parseTemplates()

	renderer := &site.Template{
		Templates: templates}

	e := echo.New()
	e.Renderer = renderer
	e.Static("/assets", "./assets/_")
	e.Use(middleware.Logger())

	e.GET("/sessions/", sessions.Index)
	e.GET("/sessions/:id", sessions.View)

	e.GET("/api/v1/sessions", sessionsAPI.Get)
	e.GET("/api/v1/sessions/:id", sessionsAPI.Get)
	e.POST("/api/v1/sessions", sessionsAPI.Post)

	e.POST("/api/v1/jobs/acquire", jobs.Acquire)

	e.POST("/api/v1/cases/:id/tries", tries.Post)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.GetApplicationConfig().Port)))

	return nil
}
