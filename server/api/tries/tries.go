package tries

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strconv"

	"github.com/graph-uk/combat-server/data/models/status"

	"github.com/graph-uk/combat-server/data/repositories"

	"github.com/graph-uk/combat-server/data/models"

	tries "github.com/graph-uk/combat-server/server/api/tries/models"

	"github.com/graph-uk/combat-server/utils"
	"github.com/labstack/echo"
)

// Post ...
func Post(c echo.Context) error {
	model := &tries.TryPostModel{}
	caseID, err := strconv.Atoi(c.Param("id"))

	casesRepo := &repositories.Cases{}
	triesRepo := &repositories.Tries{}

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	sessionCase := casesRepo.Find(caseID)

	if sessionCase == nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := c.Bind(&model); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if isTryOutFalseNegative(model.Output) {
		sessionCase.Status = status.Pending
		casesRepo.Update(sessionCase)
		return c.NoContent(http.StatusOK)
	}

	try := &models.Try{
		CaseID:     caseID,
		ExitStatus: model.ExitStatus,
		Output:     model.Output}

	tryContent, err := base64.StdEncoding.DecodeString(model.Content)

	triesRepo.Create(try, tryContent)

	return c.JSON(http.StatusOK, try)
}

func isTryOutFalseNegative(output string) bool {
	patterns := utils.GetApplicationConfig().FalseNegativePatterns
	for _, pattern := range patterns {
		r, _ := regexp.Compile(pattern)
		if r.MatchString(output) {
			return true
		}
	}
	return false
}
