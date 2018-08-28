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

func PageSessionStatusHandler(startPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// showSessionsPage(w)
		// if r.Method == "GET" {

		// 	path := strings.Split(r.URL.Path, "/")
		// 	sessionID := strings.TrimSpace(path[len(path)-1])
		// 	if strings.TrimSpace(sessionID) == "" {
		// 		//			w.Write([]byte("Session ID is not specified. Please, provide session ID like: /sessions/11203487203498"))
		// 		showSessionsPage(w)
		// 		return
		// 	}

		// 	// Get session status
		// 	Page_session, err := session.NewAssignedSession(sessionID, &t.mdb, startPath)
		// 	if err != nil {
		// 		w.Write([]byte("Error: " + err.Error() + "<br>\n"))
		// 		return
		// 	}

		// 	PS_session, err := Page_session.GetSessionPageStruct()
		// 	if err != nil {
		// 		w.Write([]byte("Error: " + err.Error()))
		// 		return
		// 	}

		// 	// Create a template.
		// 	pageBuffer := new(bytes.Buffer)

		// 	tt, err := template.ParseFiles("*views/sessions/view.html")
		// 	if err != nil {
		// 		w.Write([]byte("Error: " + err.Error() + "<br>\n"))
		// 		return
		// 	}

		// 	err = tt.Execute(pageBuffer, PS_session)
		// 	if err != nil {
		// 		w.Write([]byte("Error: " + err.Error() + "<br>\n"))
		// 		return
		// 	}

		// 	w.Write(pageBuffer.Bytes())
		// }
	}
}

// Index sessions page
func Index(c echo.Context) error {
	repo := &repositories.Sessions{}
	model := &sessions.List{
		ProjectName: configuration.ProjectName,
		Sessions:    repo.FindAll()}

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

	timestamp, _ := strconv.ParseInt(session.ID, 10, 64)
	time := time.Unix(timestamp/(1000*int64(time.Millisecond)), 0).Format("2006-01-02 15:04:05")

	cases := casesRepo.FindBySessionID(session.ID)

	model := &sessions.View{
		ProjectName: configuration.ProjectName,
		Session:     *session,
		Time:        time,
		Cases:       cases,
	}

	return c.Render(http.StatusOK, "sessions/views/view.html", model)
}
