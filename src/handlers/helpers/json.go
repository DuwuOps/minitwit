package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

func ExtractJson(c echo.Context) (map[string]any, error) {

	jsonBody := make(map[string]any)
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		fmt.Printf("json.NewDecoder returned error: %v\n", err)
		return nil, err
	}

	return jsonBody, nil
}

func GetStringValue(jsonBody map[string]any, key string) string {
	result := jsonBody[key]
	if result == nil {
		fmt.Printf("result of %v: nil\n", key)
		return ""
	}
	fmt.Printf("result.(string) of %v: %v\n", key, result.(string))
	return result.(string)
}
