package handlers

import (
	"context"
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


func getUserByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("User not found for ID: %d", userID)
		return nil, err
	}
	return user, nil
}

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

func GetCurrentUser(c echo.Context) (*models.User, error) {
	id, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("GetCurrentUser: getSessionUserID returned error: %v\n", err)
		return nil, err
	}

	user, err := userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Printf("GetCurrentUser: userRepo.GetByID returned error: %v\n", err)
	}
	return user, nil
}