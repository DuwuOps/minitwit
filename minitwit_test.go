package main

// Reference: https://echo.labstack.com/docs/testing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func register(username string, password string, password2 string, email string) echo.Context {
	if password2 != "" {
		password2 = password
	}
	if email != "" {
		email = username + "@example.com"
	}
	data := fmt.Sprintf(`{"username":%s,"password":%s, "password2":%s, "email":%s}`,
		username,
		password,
		password2,
		email,
	)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(data))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c
}

//curl -v -d "username=user1&email=user1@mail.com&password=123&password2=123" http://localhost:8000/register -L

func TestRegister(t *testing.T) {
	// c := register("user1", "default", "", "")

	// t.Log(c.Response())
	// if assert.NoError(t, Register(c)) {
	// 	assert.Equal(t, http.StatusCreated, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }

	t.Run("should return 200 status ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/register", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		Register(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
