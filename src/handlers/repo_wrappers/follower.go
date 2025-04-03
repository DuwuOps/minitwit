package repo_wrappers

import (
	"context"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"

	"github.com/labstack/echo/v4"
)

func CreateFollower(c echo.Context, whoID int, whomID int) error {
	newFollower := helpers.NewFollower(whoID, whomID)
	err := followerRepo.Create(context.Background(), newFollower)
	if err != nil {
		log.Printf("followerRepo.Create returned error: %v\n", err)
		return err
	}
	return nil
}

func DeleteFollower(c echo.Context, whoID int, whomID int) error {
	conditions := map[string]any{
		"who_id": whoID,
		"WHOM_ID": whomID,
	}
	err := followerRepo.DeleteByFields(c.Request().Context(), conditions)
	if err != nil {
		log.Printf("followerRepo.DeleteByFields returned error: %v\n", err)
		return err
	}
	return nil
}

func GetFollowerFiltered(c echo.Context, conditions map[string]any, noFollowers int) ([]models.Follower, error) {
	followers, err := followerRepo.GetFiltered(c.Request().Context(), conditions, noFollowers, "")
	if err != nil {
		log.Printf("GetFollowerFiltered: followerRepo.GetFiltered returned err %v", err)
		return nil, err
	}
	return followers, nil
}