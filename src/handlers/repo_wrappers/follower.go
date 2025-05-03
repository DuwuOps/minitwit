package repo_wrappers

import (
	"context"
	"errors"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"

	"github.com/labstack/echo/v4"
)

func CreateFollower(c echo.Context, followerID int, followingID int) error {
	if followerID == 0 || followingID == 0 {
		err := errors.New("followerID and followingID must be set")
        log.Printf("CreateFollower returned error: %v\n", err)
		return err
    }
	newFollower := helpers.NewFollower(followerID, followingID)
	err := followerRepo.Create(context.Background(), newFollower)
	if err != nil {
		log.Printf("followerRepo.Create returned error: %v\n", err)
		return err
	}
	return nil
}

func DeleteFollower(c echo.Context, followerID int, followingID int) error {
	conditions := map[string]any{
		"follower_id": followerID,
		"following_id": followingID,
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

func CountFieldInRange(c context.Context, field string, lower, upper int) (int, error) {
    count, err := followerRepo.CountRowsWhenGroupedByFieldInRange(c, field, lower, upper)
	if err != nil {
		log.Printf("❌ Repository Error: CountFieldInRange returned err %v", err)
		return 0, err
	}
	return count, nil
}
