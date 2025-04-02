package helpers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

var LATEST_PROCESSED string = "../../../latest_processed_sim_action_id.txt"

func GetLatest(c echo.Context) error {
	log.Println("🎺 User entered GetLatest via route \"/latest\"")
	
	id, err := os.ReadFile(LATEST_PROCESSED)
	if err != nil {
		log.Printf("could not read from latest_processed_sim_action_id.txt: %v\n", err)
		return err
	}

	latestProcessedCommandId, err := strconv.Atoi(string(id))
	if err != nil {
		log.Printf("latestProcessedCommandId is not an int: %v\n", err)
		return err
	}

	data := map[string]any{
		"latest": latestProcessedCommandId,
	}

	return c.JSON(http.StatusOK, data)
}

func CreateLatestFile() {
	if _, err := os.Stat(LATEST_PROCESSED); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(LATEST_PROCESSED)
		if err != nil {
			log.Printf("Could not create file. %v\n", err)
			return
		}
		os.WriteFile(LATEST_PROCESSED, []byte("0"), 0644) // If latest_processed_sim_action_id.txt does not exist, create it with an initial value.
	}
}

func UpdateLatest(c echo.Context) error {
	if _, err := os.Stat(LATEST_PROCESSED); errors.Is(err, os.ErrNotExist) {
		return err
	}

	parsedCommandId := c.FormValue("latest")

	if parsedCommandId != "" {
		return os.WriteFile(LATEST_PROCESSED, []byte(parsedCommandId), 0644)
	} else {
		return nil
	}
}
