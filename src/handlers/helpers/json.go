package helpers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

func ExtractJson(c echo.Context) (map[string]any, error) {

	jsonBody := make(map[string]any)
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		utils.LogError("json.NewDecoder returned an error", err)
		return nil, err
	}

	return jsonBody, nil
}

func GetStringValue(jsonBody map[string]any, key string) string {
	result := jsonBody[key]
	if result == nil {
		slog.Info("GetStringValue: jsonBody does not contain key", slog.Any("key", key))
		return ""
	}

	resultStr, ok := result.(string)
	if !ok {
		slog.Warn("Result is not a string!", slog.Any("key", key),
			slog.Any("value", result),
			slog.String("type", fmt.Sprintf("%T", result)))
	}

	slog.Info("GetStringValue: succesful", slog.Any("key", key), slog.Any("result", resultStr))
	return resultStr
}
