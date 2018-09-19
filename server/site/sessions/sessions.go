package sessions

import (
	"net/http"
	"strconv"
	"time"

	"github.com/graph-uk/combat-server/data/repositories"
	"github.com/graph-uk/combat-server/server/config"
	sessions "github.com/graph-uk/combat-server/server/site/sessions/models"
	"github.com/labstack/echo"
)

type handlerConfig struct {
	method  string
	route   string
	handler func([]string, http.ResponseWriter)
}

var configuration *config.Config
var handlers []*handlerConfig

func init() {
	configuration, _ = config.LoadConfig()
}

func getSessionItems() []sessions.SessionItem {
	repo := &repositories.Sessions{}

	var result []sessions.SessionItem
	items := repo.FindAll()

	for index, item := range items {
		timestamp, _ := strconv.ParseInt(item.ID, 10, 64)

		result = append(result, sessions.SessionItem{
			ID:     item.ID,
			Index:  index + 1,
			Time:   time.Unix(timestamp/(1000*int64(time.Millisecond)), 0).Format("2006-01-02 15:04:05"),
			Status: "success"})
	}

	return result
}

// Index sessions page
func Index(c echo.Context) error {
	model := &sessions.List{
		ProjectName: configuration.ProjectName,
		Sessions:    getSessionItems()}

	return c.Render(http.StatusOK, "sessions/views/index.html", model)
}

// View session
func View(c echo.Context) error {
	sessionsRepo := &repositories.Sessions{}
	casesRepo := &repositories.Cases{}
	// triesRepo := &repositories.Tries{}

	sessionID := c.Param("id")
	session := sessionsRepo.Find(sessionID)

	if session == nil {
		return c.NoContent(http.StatusNotFound)
	}

	time := session.DateCreated.Format("2006-01-02 15:04:05")

	cases := casesRepo.FindBySessionID(session.ID)

	model := &sessions.View{
		ProjectName: configuration.ProjectName,
		Session:     *session,
		Time:        time,
		Cases:       cases,
	}

	return c.Render(http.StatusOK, "sessions/views/view.html", model)
}
