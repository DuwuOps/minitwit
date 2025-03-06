package helpers

import (
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
	sess.Save(c.Request(), c.Response())
	return nil
}

func GetFlashes(c echo.Context) ([]string, error) {
	sess, err := GetSession(c)
	if err != nil {
		return []string{}, err
	}
	flashes, ok := sess.Values["Flashes"].([]string)
	if !ok {
		return []string{}, nil
	}
	sess.Values["Flashes"] = []string{}
	sess.Save(c.Request(), c.Response())
	return flashes, nil
}