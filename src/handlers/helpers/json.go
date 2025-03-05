package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

func ExtractJson(c echo.Context) (error, map[string]any) {

	jsonBody := make(map[string]any)
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		fmt.Printf("json.NewDecoder returned error: %v\n", err)
		return err, nil
	}

	return nil, jsonBody
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
