package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"log"

	"github.com/labstack/echo/v4"
)

var LATEST_PROCESSED string = "../../../latest_processed_sim_action_id.txt"

func ValidateRequest(c echo.Context) error {
	if err := UpdateLatest(c); err != nil {
		log.Printf("Error updating latest: %v", err)
		return err
	}
	if err := NotReqFromSimulator(c); err != nil {
		return err
	}
	return nil
}

func GetLatest(c echo.Context, db *sql.DB) error {
	id, err := os.ReadFile(LATEST_PROCESSED)
	if err != nil {
		fmt.Printf("could not read from latest_processed_sim_action_id.txt: %v\n", err)
		return err
	}

	latestProcessedCommandId, err := strconv.Atoi(string(id))
	if err != nil {
		fmt.Printf("latestProcessedCommandId is not an int: %v\n", err)
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
			fmt.Printf("Could not create file. %v\n", err)
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

