package handlers

import (
	"errors"
	"fmt"
	"log"
	"minitwit/src/datalayer"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"
	"net/http"
	"strings"
	"context"
	"io"
	"bytes"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var userRepo *datalayer.Repository[models.User]
var userSqliteRepo *datalayer.SqliteRepository[models.User]

func SetUserRepo(repo *datalayer.Repository[models.User], sqliteRepo *datalayer.SqliteRepository[models.User]) {
	userRepo = repo
	userSqliteRepo = sqliteRepo
}

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

		user, err := userSqliteRepo.GetByField(context.Background(), "username", username)
		if errors.Is(err, datalayer.ErrRecordNotFound) {
			errorMessage = "Invalid username"
		} else if err != nil {
			fmt.Printf("GetByField returned error: %v\n", err)
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

	if loggedIn, _ := helpers.IsUserLoggedIn(c); loggedIn {
		return c.Redirect(http.StatusFound, "/")
	}

	if err := helpers.UpdateLatest(c); err != nil {
		log.Printf("helpers.UpdateLatest returned error: %v\n", err)
		return err
	}

	if c.Request().Method == http.MethodPost {
		username, email, password, password2, pwd, err := extractRegisterInput(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error_msg": "Invalid request payload"})
		}

		if errorMessage := validateRegisterInput(username, email, password, password2); errorMessage != "" {
			return renderRegisterError(c, errorMessage)
		}

		if existingUser, _ := userSqliteRepo.GetByField(context.Background(), "username", username); existingUser != nil {
			return renderRegisterError(c, "The username is already taken")
		}

		if err := createUser(username, email, password); err != nil {
			log.Printf("CreateUser returned error: %v\n", err)
			return err
		}

		if pwd != "" {
			return c.String(http.StatusNoContent, "")
		}

		helpers.AddFlash(c, "You were successfully registered and can login now")
		return c.Redirect(http.StatusFound, "/login")
	}

	return renderRegisterPage(c, "")
}


func renderRegisterError(c echo.Context, errorMessage string) error {
	flashes, _ := helpers.GetFlashes(c)
	data := map[string]any{
		"Error":   errorMessage,
		"Flashes": flashes,
	}
	return c.Render(http.StatusOK, "register.html", data)
}

func renderRegisterPage(c echo.Context, errorMessage string) error {
	flashes, _ := helpers.GetFlashes(c)
	data := map[string]any{
		"Error":   errorMessage,
		"Flashes": flashes,
	}
	return c.Render(http.StatusOK, "register.html", data)
}

func extractRegisterInput(c echo.Context) (string, string, string, string, string, error) {
	bodyBytes, _ := io.ReadAll(c.Request().Body)
	log.Println("📩 Raw Request Body:", string(bodyBytes))

	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	contentType := c.Request().Header.Get("Content-Type")

	var username, email, password, password2, pwd string

	if contentType == "application/json" {
		payload := make(map[string]string)
		if err := json.Unmarshal(bodyBytes, &payload); err == nil {
			username = payload["username"]
			email = payload["email"]
			password = payload["password"]
			password2 = payload["password2"]
			pwd = payload["pwd"]
		} else {
			log.Printf("❌ JSON Decoding Error: %v\n", err)
		}
	} else {
		log.Println("✅ Extracted Form Payload (Fallback)")
		username = c.FormValue("username")
		email = c.FormValue("email")
		password = c.FormValue("password")
		password2 = c.FormValue("password2")
		pwd = c.FormValue("pwd")
	}

	if password == "" {
		password = pwd
		password2 = pwd
	}

	return username, email, password, password2, pwd, nil
}



func validateRegisterInput(username, email, password, password2 string) string {
	switch {
	case username == "":
		return "You have to enter a username"
	case email == "" || !strings.Contains(email, "@"):
		return "You have to enter a valid email address"
	case password == "":
		return "You have to enter a password"
	case password != password2:
		return "The two passwords do not match"
	default:
		return ""
	}
}

func createUser(username, email, password string) error {
    log.Printf("🔹 Creating user: Username=%s, Email=%s", username, email)

    hash, err := generatePasswordHash(password)
    if err != nil {
        return fmt.Errorf("error hashing password: %v", err)
    }

    newUser := &models.User{
        Username: username,
        Email:    email,
        PwHash:   hash,
    }

    err = userSqliteRepo.Create(context.Background(), newUser)
    if err != nil {
        log.Printf("❌ Error inserting user into DB: %v", err)
        return err
    }

	err = userRepo.Create(context.Background(), newUser)
    if err != nil {
        log.Printf("❌ Error inserting user into Postgres DB: %v", err)
        return err
    }

    log.Println("✅ User successfully inserted into DB")
    return nil
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
