package sessions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/graph-uk/combat-server/data/repositories"
	sessions "github.com/graph-uk/combat-server/server/site/sessions/models"
	"github.com/graph-uk/combat-server/utils"
	"github.com/labstack/echo"
)

type handlerConfig struct {
	method  string
	route   string
	handler func([]string, http.ResponseWriter)
}

var handlers []*handlerConfig

func getSessionItems() []sessions.SessionItem {
	repo := &repositories.Sessions{}

	var result []sessions.SessionItem
	items := repo.FindAll()

	for index, item := range items {

		result = append(result, sessions.SessionItem{
			ID:     item.ID,
			Index:  index + 1,
			Time:   item.DateCreated.Format("2006-01-02 15:04:05"),
			Status: strings.ToLower(item.Status.String())})
	}

	return result
}

// Index sessions page
func Index(c echo.Context) error {
	model := &sessions.List{
		ProjectName: utils.GetApplicationConfig().ProjectName,
		Sessions:    getSessionItems()}

	return c.Render(http.StatusOK, "sessions/views/index.html", model)
}

func getCasesJSON(sessionID string) string {
	model := map[string]sessions.CaseItem{}
	triesRepo := &repositories.Tries{}
	casesRepo := &repositories.Cases{}

	cases := casesRepo.FindBySessionID(sessionID)

	for _, sessionCase := range cases {
		var tries []sessions.TryItem
		id := fmt.Sprintf("case%d", sessionCase.ID)
		caseTries := triesRepo.FindByCaseID(sessionCase.ID)

		//
		for _, try := range caseTries {
			var steps []sessions.TryStepItem
			rawSteps := triesRepo.FindTrySteps(try.ID)

			for _, step := range rawSteps {
				stepUrlBuff, _ := ioutil.ReadFile(fmt.Sprintf("./_data/tries/%d/_/out/%s.txt", try.ID, step))

				steps = append(steps, sessions.TryStepItem{
					Image:  fmt.Sprintf("/tries/%d/%s.png", try.ID, step),
					Source: fmt.Sprintf("/tries/%d/%s.html", try.ID, step),
					URL:    string(stepUrlBuff),
				})
			}

			tries = append(tries, sessions.TryItem{
				Output: try.Output,
				Steps:  steps})
		}

		var steps []sessions.TryStepItem
		lastSuccessfulRun := triesRepo.FindLastSuccessfulTry(sessionCase.ID)
		rawSteps := triesRepo.FindTrySteps(lastSuccessfulRun.ID)
		for _, step := range rawSteps {
			stepUrlBuff, _ := ioutil.ReadFile(fmt.Sprintf("./_data/tries/%d/_/out/%s.txt", lastSuccessfulRun.ID, step))

			steps = append(steps, sessions.TryStepItem{
				Image:  fmt.Sprintf("/tries/%d/%s.png", lastSuccessfulRun.ID, step),
				Source: fmt.Sprintf("/tries/%d/%s.html", lastSuccessfulRun.ID, step),
				URL:    string(stepUrlBuff),
			})
		}

		lastSuccessfulRunTryItem := sessions.TryItem{}
		lastSuccessfulRunTryItem.Output = lastSuccessfulRun.Output
		lastSuccessfulRunTryItem.Steps = steps

		model[id] = sessions.CaseItem{
			Status:            strings.ToLower(sessionCase.Status.String()),
			Title:             sessionCase.Title,
			Tries:             tries,
			LastSuccessfulRun: lastSuccessfulRunTryItem,
		}
	}

	result, _ := json.Marshal(model)
	return string(result)
}

// View session
func View(c echo.Context) error {
	sessionsRepo := &repositories.Sessions{}

	sessionID := c.Param("id")
	session := sessionsRepo.Find(sessionID)

	if session == nil {
		return c.NoContent(http.StatusNotFound)
	}

	title := session.DateCreated.Format("2006-01-02 15:04:05")

	model := &sessions.View{
		ProjectName: utils.GetApplicationConfig().ProjectName,
		Title:       title,
		Cases:       getCasesJSON(session.ID)}

	return c.Render(http.StatusOK, "sessions/views/view.html", model)
}
