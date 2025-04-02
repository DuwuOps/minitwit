package helpers

import (
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
)

func ExtractJson(c echo.Context) (map[string]any, error) {

	jsonBody := make(map[string]any)
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Printf("json.NewDecoder returned error: %v\n", err)
		return nil, err
	}

	return jsonBody, nil
}

func GetStringValue(jsonBody map[string]any, key string) string {
	result := jsonBody[key]
	if result == nil {
		log.Printf("result of %v: nil\n", key)
		return ""
	}
	log.Printf("result.(string) of %v: %v\n", key, result.(string))
	return result.(string)
}
