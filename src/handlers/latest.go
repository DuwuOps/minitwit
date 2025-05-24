package handlers

import (
	"log/slog"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)


func GetLatest(c echo.Context) error {
	
	utils.LogRouteStart(c, "GetLatest", "/latest")

	latestProcessedCommandStr, err := repo_wrappers.GetLatest(c)
	if err != nil {
		return err
	}

	latestProcessedCommandId := int64(0)

	if len(latestProcessedCommandStr) == 0 {
		slog.InfoContext(c.Request().Context(), "latestProcessedCommandId not found.")
	} else {
		latestProcessedCommandId = latestProcessedCommandStr[0].LatestProcessedID
	}

	data := map[string]any{
		"latest": latestProcessedCommandId,
	}

	return c.JSON(http.StatusOK, data)
}
