package handlers

import (
	"errors"
	"minitwit/src/models"
)


var ErrRecordNotFound = errors.New("record not found")

func newUser(username string, email string, hash string) *models.User {
	return &models.User{
		Username: username,
		Email: email,
		PwHash: hash,
	}
}