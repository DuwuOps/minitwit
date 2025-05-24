package handlers

import (
	"context"
	"errors"
	"minitwit/src/datalayer"
	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	utils.LogRouteStart(c, "Login", "/login")
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		return c.Redirect(http.StatusFound, "/")
	}

	var errorMessage string
	if c.Request().Method == http.MethodPost {
		username := c.FormValue("username")
		password := c.FormValue("password")


		user, err := repo_wrappers.GetUserByUsername(context.Background(), username)

		if errors.Is(err, datalayer.ErrRecordNotFound) {
			errorMessage = "Invalid username"
		} else if err != nil {
			utils.LogError("Db.QueryRow returned an error", err)
			return err
		} else {
			if !checkPasswordHash(user.PwHash, password) {
				errorMessage = "Invalid password"
			} else {
				err = helpers.AddFlash(c, "You were logged in")
				if err != nil {
					utils.LogError("addFlash returned an error", err)
				}
				helpers.SetSessionUserID(c, user.UserID)
				return c.Redirect(http.StatusFound, "/")
			}
		}
	}

	flashes, _ := helpers.GetFlashes(c)

	data := map[string]any{
		"Error":   errorMessage,
		"Flashes": flashes,
	}
	return c.Render(http.StatusOK, "login.html", data)
}

func Register(c echo.Context) error {
	utils.LogRouteStart(c, "Register", "/register")
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		return c.Redirect(http.StatusFound, "/")
	}

	err := repo_wrappers.UpdateLatest(c)
	if err != nil {
		utils.LogError("helpers.UpdateLatest returned an error", err)
		return err
	}

	var errorMessage string
	if c.Request().Method == http.MethodPost {
		payload, err := helpers.ExtractJson(c)
		if err != nil {
			utils.LogErrorContext(c.Request().Context(), "Register: ExtractJson returned an error", err)
		}

		var username string
		var email string
		var pwd string
		var password string
		var password2 string

		if payload != nil {
			username = helpers.GetStringValue(payload, "username")
			email = helpers.GetStringValue(payload, "email")
			pwd = helpers.GetStringValue(payload, "pwd")
			password = helpers.GetStringValue(payload, "password")
			password2 = helpers.GetStringValue(payload, "password2")
		} else {
			username = c.FormValue("username")
			email = c.FormValue("email")
			pwd = c.FormValue("pwd")
			password = c.FormValue("password")
			password2 = c.FormValue("password2")
		}

		if password == "" {
			password = pwd
			password2 = pwd
		}

		switch {
		case username == "":
			errorMessage = "You have to enter a username"
		case email == "" || !strings.Contains(email, "@"):
			errorMessage = "You have to enter a valid email address"
		case password == "":
			errorMessage = "You have to enter a password"
		case password != password2:
			errorMessage = "The two passwords do not match"
		default:
			existingUser, _ := repo_wrappers.GetUserByUsername(context.Background(), username)
			if existingUser != nil {
				errorMessage = "The username is already taken"
			} else {
				hash, err := generatePasswordHash(password)
				if err != nil {
					utils.LogError("generatePasswordHash returned an error", err)
					return err
				}
				
				_ = repo_wrappers.CreateUser(username, email, hash)
				

				if pwd == "" {
					err = helpers.AddFlash(c, "You were successfully registered and can login now")
					if err != nil {
						utils.LogError("helpers.addFlash returned an error", err)
					}
					return c.Redirect(http.StatusFound, "/login")
				}
			}
		}
		if pwd != "" {
			if errorMessage != "" {
				data := map[string]any{
					"error_msg": errorMessage,
				}
				return c.JSON(http.StatusBadRequest, data)
			}
			return c.String(http.StatusNoContent, "")
		}
	}

	flashes, _ := helpers.GetFlashes(c)

	data := map[string]any{
		"Error":   errorMessage,
		"Flashes": flashes,
	}
	return c.Render(http.StatusOK, "register.html", data)
}

func Logout(c echo.Context) error {
	utils.LogRouteStart(c, "Logout", "/logout")
	err := helpers.ClearSessionUserID(c)
	if err != nil {
		utils.LogErrorEchoContext(c, "ClearSessionUserID returned an error", err)
	}
	err = helpers.AddFlash(c, "You were logged out")
	if err != nil {
		utils.LogError("helpers.addFlash returned an error", err)
	}
	return c.Redirect(http.StatusFound, "/public")
}

func checkPasswordHash(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func generatePasswordHash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}
