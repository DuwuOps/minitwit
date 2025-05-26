package helpers

import (
	"log/slog"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

func AddFlash(c echo.Context, message string) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}
	flashes, ok := sess.Values["Flashes"].([]string)
	if !ok {
		flashes = []string{}
	}
	flashes = append(flashes, message)
	sess.Values["Flashes"] = flashes
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		utils.LogErrorEchoContext(c, "Session.Save returned an error", err)
	}
	return nil
}

func GetFlashes(c echo.Context) ([]string, error) {
	sess, err := GetSession(c)
	if err != nil {
		return []string{}, err
	}
	flashes, ok := sess.Values["Flashes"].([]string)
	if !ok {
		utils.InfoEchoContext(c, "0 Flashes found")
		return []string{}, nil
	}
	sess.Values["Flashes"] = []string{}
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		utils.LogErrorEchoContext(c, "Session.Save returned an error", err)
	}
	slog.InfoContext(c.Request().Context(), "Flashes found", slog.Any("flashed_found", len(flashes)), slog.Any("flashes", flashes))
	return flashes, nil
}
