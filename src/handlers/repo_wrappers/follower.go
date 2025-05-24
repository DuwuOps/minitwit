package repo_wrappers

import (
	"context"
	"errors"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

func CreateFollower(c echo.Context, followerID int, followingID int) error {
	if followerID == 0 || followingID == 0 {
		err := errors.New("followerID and followingID must be set")
        utils.LogErrorEchoContext(c, "CreateFollower returned an error", err)
		return err
    }
	newFollower := helpers.NewFollower(followerID, followingID)
	err := followerRepo.Create(context.Background(), newFollower)
	if err != nil {
		utils.LogErrorEchoContext(c, "followerRepo.Create returned an error", err)
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
		utils.LogErrorEchoContext(c, "followerRepo.DeleteByFields returned an error", err)
		return err
	}
	return nil
}

func GetFollowerFiltered(c echo.Context, conditions map[string]any, noFollowers int) ([]models.Follower, error) {
	followers, err := followerRepo.GetFiltered(c.Request().Context(), conditions, noFollowers, "")
	if err != nil {
		utils.LogErrorEchoContext(c, "GetFollowerFiltered: followerRepo.GetFiltered returned an error", err)
		return nil, err
	}
	return followers, nil
}

func CountFieldInRange(c context.Context, field string, lower, upper int) (int, error) {
    count, err := followerRepo.CountRowsWhenGroupedByFieldInRange(c, field, lower, upper)
	if err != nil {
		utils.LogErrorContext(c, "❌ Repository Error: CountFieldInRange returned an error", err)
		return 0, err
	}
	return count, nil
}
