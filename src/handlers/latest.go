package handlers

import (
	"log"
	"minitwit/src/handlers/repo_wrappers"
	"net/http"

	"github.com/labstack/echo/v4"
)


func GetLatest(c echo.Context) error {
	
	log.Println("🎺 User entered GetLatest via route \"/latest\"")

	latestProcessedCommandStr, err := repo_wrappers.GetLatest(c)
	if err != nil {
		return err
	}

	latestProcessedCommandId := 0

	if len(latestProcessedCommandStr) == 0 {
		log.Printf("latestProcessedCommandId not found.\n")
	} else {
		latestProcessedCommandId = latestProcessedCommandStr[0].LatestProcessedID
	}

	data := map[string]any{
		"latest": latestProcessedCommandId,
	}

	return c.JSON(http.StatusOK, data)
}
