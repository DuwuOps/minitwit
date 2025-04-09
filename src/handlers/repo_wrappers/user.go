package repo_wrappers

import (
	"context"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/metrics"
	"minitwit/src/models"

	"github.com/labstack/echo/v4"
)

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := userRepo.GetByField(ctx, "username", username)
	if err != nil {
		log.Printf("User not found: %s", username)
		return nil, err
	}
	return user, nil
}

func GetUserByID(c echo.Context, id int) (*models.User, error) {
	user, err := userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Printf(`User not found: %d`, id)
		return nil, err
	}
	return user, nil
}

func IsFollowingUser(c echo.Context, profileUserID int) bool {
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

func GetCurrentUser(c echo.Context) (*models.User, error) {
	id, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("GetCurrentUser: getSessionUserID returned error: %v\n", err)
		return nil, err
	}

	user, err := userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Printf("GetCurrentUser: userRepo.GetByID returned error: %v\n", err)
		return nil, err
	}
	return user, nil
}

func CreateUser(username string, email string, hash string) error {
	newUser := helpers.NewUser(username, email, hash)
	err := userRepo.Create(context.Background(), newUser)
	if err != nil {
		log.Printf("userRepo.Create returned error: %v\n", err)
		return err
	}
	metrics.NewUsers.Inc()
	return nil
}

func GetUserMap(c echo.Context, userIDs []int) map[string]*models.User {
	return nil
}