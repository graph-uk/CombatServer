package configs

import (
	"net/http"
	"time"

	"malibu-server/data/repositories"
	"malibu-server/server/api/configs/models"
	"malibu-server/utils"

	"github.com/labstack/echo"
)

// Get current session status
func Get(c echo.Context) error {
	configsRepo := &repositories.Configs{}
	dbConfig := configsRepo.Find()

	appConfig := utils.GetApplicationConfig()

	result := &configs.ConfigModel{
		NotificationEnabled: dbConfig.NotificationEnabled,
		MuteTimestamp:       dbConfig.MuteTimestamp,
		MuteDurationMinutes: appConfig.NotificationMuteDurationMinutes,
	}

	return c.JSON(http.StatusOK, result)
}

func Put(c echo.Context) error {
	model := &configs.ConfigPutModel{}
	configsRepo := &repositories.Configs{}

	if err := c.Bind(&model); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dbConfig := configsRepo.Find()
	dbConfig.NotificationEnabled = model.NotificationEnabled
	dbConfig.MuteTimestamp = time.Now()

	configsRepo.Update(dbConfig)

	return c.JSON(http.StatusCreated, dbConfig)
}
