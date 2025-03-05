package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

var LATEST_PROCESSED string = "../../../latest_processed_sim_action_id.txt"

func RowsToMapList(rows *sql.Rows) ([]map[string]any, error) {
	var result []map[string]any
	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s (any's) to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Printf("rows.Scan returned error: %v\n", err)
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		result = append(result, m)
	}

	return result, nil
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

func UpdateLatest(c echo.Context) {
	if _, err := os.Stat(LATEST_PROCESSED); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(LATEST_PROCESSED)
		if err != nil {
			fmt.Printf("Could not create file. %v\n", err)
			return
		}
	}
	parsedCommandId := c.FormValue("latest")

	if parsedCommandId != "" {
		os.WriteFile(LATEST_PROCESSED, []byte(parsedCommandId), 0644)
	}
}
