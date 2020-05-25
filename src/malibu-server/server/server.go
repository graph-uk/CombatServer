package server

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"malibu-server/server/api/configs"
	"malibu-server/server/api/jobs"
	sessionsAPI "malibu-server/server/api/sessions"
	"malibu-server/server/api/tries"
	"malibu-server/server/site"
	"malibu-server/server/site/sessions"
	"malibu-server/utils"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// MalibuServer ...
type MalibuServer struct {
	startPath string
}

func parseTemplates() (*template.Template, error) {
	root := template.New("")
	tempBox := packr.NewBox("site")

	err := tempBox.Walk(func(path string, file packr.File) error {
		path = strings.Replace(path, `\`, `/`, -1)
		if strings.HasSuffix(path, ".html") {
			b := file.String()

			t := root.New(path)
			t, e2 := t.Parse(b)
			if e2 != nil {
				return e2
			}

		}
		return nil
	})
	return root, err
}

// Start web server
func (t *MalibuServer) Start(config *utils.Config) error {

	go TimeoutWatcher(config)

	templates, _ := parseTemplates()

	renderer := &site.Template{
		Templates: templates}

	e := echo.New()

	e.Pre(middleware.Rewrite(map[string]string{
		"/tries/*/*": "/tries/$1/_/out/$2",
	}))

	assetsBox := packr.NewBox("../assets/_")

	assetHandler := http.FileServer(assetsBox)

	e.Renderer = renderer
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))
	e.Static("/tries", "./_data/tries")
	e.Static("/tries-succ", "./_data/tries-succ")

	//e.Use(middleware.Logger())

	e.GET("/sessions/", sessions.Index)
	e.GET("/", sessions.Index)
	e.GET("/sessions/:id", sessions.View)

	e.GET("/api/v1/sessions", sessionsAPI.Get)
	e.GET("/api/v1/sessions/:id", sessionsAPI.Get)
	e.POST("/api/v1/sessions", sessionsAPI.Post)

	e.POST("/api/v1/jobs/acquire", jobs.Acquire)

	e.POST("/api/v1/cases/:id/tries", tries.Post)

	e.GET("/api/v1/config", configs.Get)
	e.PUT("/api/v1/config", configs.Put)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(utils.GetApplicationConfig().Port)))

	return nil
}
