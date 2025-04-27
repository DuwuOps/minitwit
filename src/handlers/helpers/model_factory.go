package helpers

import (
	"minitwit/src/models"
	"time"
)

func NewMessage(authorID int, text string) *models.Message {
	return &models.Message{
		AuthorID: authorID,
		Text:     text,
		PubDate:  time.Now().Unix(),
		Flagged:  0,
	}
}

func NewUser(username string, email string, hash string) *models.User {
	return &models.User{
		Username: username,
		Email: email,
		PwHash: hash,
	}
}

func NewFollower(followerID int, followingID int) *models.Follower {
	return &models.Follower{
		FollowerID: followerID,
		FollowingID: followingID,
	}
}