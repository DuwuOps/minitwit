package repo_wrappers

import (
	"minitwit/src/models"
	"minitwit/src/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetLatest(c echo.Context) ([]models.LatestProcessed, error) {
	latestProcessed, err := latestProcessedRepo.GetFiltered(c.Request().Context(), nil, 0, "")
	if err != nil {
		utils.LogErrorEchoContext(c, "could not read latest_processed from database", err)
		return nil, err
	}

	return latestProcessed, nil
}

func UpdateLatest(c echo.Context) error {
	parsedCommandId := c.FormValue("latest")

	if parsedCommandId != "" {
		parsedCommandId, err := strconv.Atoi(parsedCommandId)
		if err != nil {
			utils.LogErrorEchoContext(c, "parsedCommandId is not an int", err)
			return err
		}

		updates := map[string]any{
			"latest_processed_id": parsedCommandId,
		}

		err = latestProcessedRepo.SetAllFields(c.Request().Context(), updates)
		if err != nil {
			utils.LogErrorEchoContext(c, "could not update latest_processed_id in database", err)
		}
		return err
	} else {
		return nil
	}
}
