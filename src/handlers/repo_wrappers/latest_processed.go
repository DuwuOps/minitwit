package repo_wrappers

import (
	"log"
	"minitwit/src/models"
	"strconv"

	"github.com/labstack/echo/v4"
)


func GetLatest(c echo.Context) ([]models.LatestProcessed, error) {
	
	latestProcessed, err := latestProcessedRepo.GetFiltered(c.Request().Context(), nil, 0, "")
	if err != nil {
		log.Printf("could not read latest_processed from database: %v\n", err)
		return nil, err
	}

	return latestProcessed, nil
}

func UpdateLatest(c echo.Context) error {
	
	parsedCommandId := c.FormValue("latest")

	if parsedCommandId != "" {
		
		parsedCommandId, err := strconv.Atoi(string(parsedCommandId))
		if err != nil {
			log.Printf("parsedCommandId is not an int: %v\n", err)
			return err
		}

		updates := map[string]any{
			"latest_processed_id": parsedCommandId,
		}

		err = latestProcessedRepo.SetAllFields(c.Request().Context(), updates)
		if err != nil {
			log.Printf("could not update latest_processed_id in database: %v\n", err)
		}
		return err
	
	} else {
		return nil
	}
}
