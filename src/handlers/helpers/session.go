package helpers

import (
	"minitwit/src/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SetSessionUserID(c echo.Context, userID int) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}
	sess.Values["user_id"] = userID
	sess.Save(c.Request(), c.Response())
	return nil
}

func GetSessionUserID(c echo.Context) (int, error) {
	sess, err := GetSession(c)
	if err != nil {
		return 0, err
	}
	id, _ := sess.Values["user_id"].(int)
	return id, nil
}

func ClearSessionUserID(c echo.Context) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}
	delete(sess.Values, "user_id")
	sess.Save(c.Request(), c.Response())
	return nil
}

func IsUserLoggedIn(c echo.Context) (bool, error) {
	sess, err := GetSession(c)
	_, ok := sess.Values["user_id"].(int)
	return ok, err
}

func GetSession(c echo.Context) (*sessions.Session, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		utils.LogError("session.Get returned an error", err)
		return nil, err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return sess, nil
}
