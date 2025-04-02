package helpers

import (
	"log"

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
		log.Println("0 Flashes found")
		return []string{}, nil
	}
	sess.Values["Flashes"] = []string{}
	sess.Save(c.Request(), c.Response())
	log.Printf("%v Flashes found: %v\n", len(flashes), flashes)
	return flashes, nil
}