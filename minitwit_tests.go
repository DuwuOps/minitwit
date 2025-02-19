package main

// Reference: https://echo.labstack.com/docs/testing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func register(username string, password string, password2 string, email string) {
	if password2 != "" {
		password2 = password
	}
	if email != "" {
		email = username + "@example.com"
	}
}

//curl -v -d "username=user1&email=user1@mail.com&password=123&password2=123" http://localhost:8000/register -L

func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()

	userJSON := `{"name":"user1","email":"user1@example.com"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.SetPath("/register")

	c.SetParamNames("email")
	c.SetParamValues("jon@labstack.com")

	// Assertions
	if assert.NoError(t, Register(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}
