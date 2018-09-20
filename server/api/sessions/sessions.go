package sessions

import (
	"encoding/base64"
	"net/http"

	"github.com/graph-uk/combat-server/data/models/status"

	"github.com/graph-uk/combat-server/data/models"

	"github.com/graph-uk/combat-server/data/repositories"
	"github.com/graph-uk/combat-server/server/api/sessions/models"
	"github.com/labstack/echo"
)

// Get current session status
func Get(c echo.Context) error {
	var session *models.Session
	var casesCount, casesProcessedCount int
	var casesFailed []string

	casesRepo := &repositories.Cases{}
	sessionsRepo := &repositories.Sessions{}

	sessionID := c.Param("id")

	if sessionID != "" {
		session = sessionsRepo.Find(sessionID)
	} else {
		session = sessionsRepo.FindLast()
	}

	if session == nil {
		return c.NoContent(http.StatusNotFound)
	}

	cases := casesRepo.FindBySessionID(session.ID)

	for _, sessionCase := range cases {
		casesCount++
		if sessionCase.Status != status.Pending && sessionCase.Status != status.Processing {
			casesProcessedCount++
		}
		if sessionCase.Status == status.Failed {
			casesFailed = append(casesFailed, sessionCase.Title)
		}
	}

	result := &sessions.SessionModel{
		ID:                  session.ID,
		SessionError:        session.Error,
		Status:              session.Status.String(),
		CasesCount:          casesCount,
		CasesProcessedCount: casesProcessedCount,
		CasesFailed:         casesFailed}

	return c.JSON(http.StatusOK, result)
}

// Post creates new session
func Post(c echo.Context) error {
	model := &sessions.SessionPostModel{}
	casesRepo := &repositories.Cases{}
	sessionsRepo := &repositories.Sessions{}

	if err := c.Bind(&model); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	sessionContent, err := base64.StdEncoding.DecodeString(model.Content)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid session content")
	}

	casesRepo.StopCurrentCases()
	session := sessionsRepo.Create(model.Arguments, sessionContent)

	return c.JSON(http.StatusCreated, session)
}
