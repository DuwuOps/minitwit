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

func RowsToMapList(rows *sql.Rows) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
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
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		result = append(result, m)
	}

	return result, nil
}

func GetLatest(c echo.Context, db *sql.DB) error {
	id, err := os.ReadFile("./latest_processed_sim_action_id.txt")
	if err != nil {
		fmt.Printf("could not read from ./latest_processed_sim_action_id.txt: %v\n", err)
		return err
	}

	latestProcessedCommandId, err := strconv.Atoi(string(id))
	if err != nil {
		fmt.Printf("latestProcessedCommandId is not an int: %v\n", err)
		return err
	}

	data := map[string]interface{}{
		"latest": latestProcessedCommandId,
	}

	return c.JSON(http.StatusOK, data)
}

func UpdateLatest(c echo.Context) {
	if _, err := os.Stat("./latest_processed_sim_action_id.txt"); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create("./latest_processed_sim_action_id.txt")
		if err != nil {
			fmt.Printf("Could not create file. %v\n", err)
			return
		}
	}
	parsedCommandId := c.FormValue("latest")

	if parsedCommandId != "" {
		os.WriteFile("./latest_processed_sim_action_id.txt", []byte(parsedCommandId), 0644)
	}
}
