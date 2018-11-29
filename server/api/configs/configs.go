package configs

import (
	"net/http"

	"github.com/graph-uk/combat-server/data/repositories"
	"github.com/graph-uk/combat-server/server/api/configs/models"
	"github.com/labstack/echo"
)

// Get current session status
func Get(c echo.Context) error {
	configsRepo := &repositories.Configs{}
	dbConfig := configsRepo.Find()

	result := &configs.ConfigModel{
		NotificationEnabled: dbConfig.NotificationEnabled,
		MuteTimestamp:       dbConfig.MuteTimestamp,
	}

	return c.JSONPretty(http.StatusOK, result, ` `)
}

func Put(c echo.Context) error {
	model := &configs.ConfigPutModel{}
	configsRepo := &repositories.Configs{}

	if err := c.Bind(&model); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dbConfig := configsRepo.Find()
	dbConfig.NotificationEnabled = model.NotificationEnabled
	dbConfig.MuteTimestamp = model.MuteTimestamp

	configsRepo.Update(dbConfig)

	return c.JSON(http.StatusCreated, dbConfig)
}
