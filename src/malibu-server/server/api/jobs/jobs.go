package jobs

import (
	"encoding/base64"
	"net/http"

	"malibu-server/data/repositories"
	jobs "malibu-server/server/api/jobs/models"

	"github.com/labstack/echo"
)

func Acquire(c echo.Context) error {
	repo := &repositories.Cases{}
	sessionFSRepo := &repositories.SessionsFS{}

	sessionCase := repo.AcquireFreeJob()

	if sessionCase == nil {
		return c.NoContent(http.StatusNotFound)
	}

	sessionContent := sessionFSRepo.FindSessionContent(sessionCase.SessionID)

	if sessionContent == nil {
		return c.NoContent(http.StatusNotFound)
	}

	sessionContentEncoded := base64.StdEncoding.EncodeToString(sessionContent)

	model := jobs.Job{
		CaseID:      sessionCase.ID,
		CommandLine: sessionCase.CommandLine,
		Content:     sessionContentEncoded}

	return c.JSON(http.StatusOK, model)
}
