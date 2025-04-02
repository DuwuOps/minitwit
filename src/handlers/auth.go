package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/datalayer"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	log.Println("User entered Login via route \"/login\"")
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		return c.Redirect(http.StatusFound, "/")
	}

	var errorMessage string
	if c.Request().Method == http.MethodPost {
		username := c.FormValue("username")
		password := c.FormValue("password")


		user, err := userRepo.GetByField(c.Request().Context(), "username", username)

		if errors.Is(err, datalayer.ErrRecordNotFound) {
			errorMessage = "Invalid username"
		} else if err != nil {
			fmt.Printf("Db.QueryRow returned error: %v\n", err)
			return err
		} else {
			if !checkPasswordHash(user.PwHash, password) {
				errorMessage = "Invalid password"
			} else {
				helpers.AddFlash(c, "You were logged in")
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
	log.Printf("User entered Register via route \"/register\" and HTTP method %v", c.Request().Method)
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		return c.Redirect(http.StatusFound, "/")
	}

	err := helpers.UpdateLatest(c)
	if err != nil {
		fmt.Printf("helpers.UpdateLatest returned error: %v\n", err)
		return err
	}


	var errorMessage string
	if c.Request().Method == http.MethodPost {
		payload, err := helpers.ExtractJson(c)

		var username string
		var email string
		var pwd string
		var password string
		var password2 string

		if err == nil {
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
			existingUser, _ := userRepo.GetByField(context.Background(), "username", username)
			if existingUser != nil {
				errorMessage = "The username is already taken"
			} else {
				hash, err := generatePasswordHash(password)
				if err != nil {
					fmt.Printf("generatePasswordHash returned error: %v\n", err)
					return err
				}

				err = userRepo.Create(c.Request().Context(), newUser(username, email, hash))
				if err != nil {
					fmt.Printf("userRepo.Create returned error: %v\n", err)
					return err
				}

				if pwd == "" {
					helpers.AddFlash(c, "You were successfully registered and can login now")
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
	helpers.ClearSessionUserID(c)
	helpers.AddFlash(c, "You were logged out")
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
