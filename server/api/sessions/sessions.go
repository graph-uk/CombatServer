package sessions

import (
	"encoding/base64"
	"net/http"

	"github.com/graph-uk/combat-server/data/repositories"
	"github.com/graph-uk/combat-server/server/api/sessions/models"
	"github.com/labstack/echo"
)

// Post creates new session
func Post(c echo.Context) error {
	model := &sessions.SessionPostModel{}
	repo := &repositories.Cases{}
	sessionsRepo := &repositories.Sessions{}

	if err := c.Bind(&model); err != nil {
		return err
	}

	sessionContent, err := base64.StdEncoding.DecodeString(model.Content)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid session content")
	}

	repo.StopCurrentCases()
	session := sessionsRepo.Create(model.Arguments, sessionContent)

	return c.JSON(http.StatusCreated, session)
}
