package handlers

import (
	"errors"
	"minitwit/src/models"
	"time"
)


var ErrRecordNotFound = errors.New("record not found")

func newMessage(authorID int, text string) *models.Message {
	return &models.Message{
		AuthorID: authorID,
		Text:     text,
		PubDate:  time.Now().Unix(),
		Flagged:  0,
	}
}

func newUser(username string, email string, hash string) *models.User {
	return &models.User{
		Username: username,
		Email: email,
		PwHash: hash,
	}
}