package jobs

import (
	"net/http"

	"github.com/graph-uk/combat-server/data/repositories"
	jobs "github.com/graph-uk/combat-server/server/api/jobs/models"
	"github.com/labstack/echo"
)

func Acquire(c echo.Context) error {
	repo := &repositories.Cases{}
	sessionRepo := &repositories.Sessions{}

	sessionCase := repo.AcquireFreeJob()

	if sessionCase == nil {
		return c.NoContent(http.StatusNotFound)
	}

	sessionContent := sessionRepo.FindSessionContent(sessionCase.SessionID)

	if sessionContent == "" {
		return c.NoContent(http.StatusNotFound)
	}

	model := jobs.Job{
		CaseID:      sessionCase.ID,
		CommandLine: sessionCase.CommandLine,
		Content:     sessionContent}

	return c.JSON(http.StatusOK, model)
}
