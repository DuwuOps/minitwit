package handlers

import (
	"context"
	"errors"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"
	"time"

	"github.com/labstack/echo/v4"
)

func getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := userRepo.GetByField(ctx, "username", username)
	if err != nil {
		log.Printf("User not found: %s", username)
		return nil, err
	}
	return user, nil
}

func isFollowingUser(c echo.Context, profileUserID int) bool {
	sessionUserID, err := helpers.GetSessionUserID(c)
	if err != nil || sessionUserID == 0 {
		return false
	}

	conditions := map[string]any{
		"who_id":  sessionUserID,
		"whom_id": profileUserID,
	}
	followers, err := followerRepo.GetFiltered(c.Request().Context(), conditions, 1, "")

	return err == nil && len(followers) > 0
}

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

func newFollower(whoID int, whomID int) *models.Follower {
	return &models.Follower{
		WhoID: whoID,
		WhomID: whomID,
	}
}