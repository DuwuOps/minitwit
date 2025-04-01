package handlers

import "minitwit/src/models"

func newUser(username string, email string, hash string) *models.User {
	return &models.User{
		Username: username,
		Email: email,
		PwHash: hash,
	}
}