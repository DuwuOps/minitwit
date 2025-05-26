package repo_wrappers

import (
	"context"
	"log/slog"
	"minitwit/src/handlers/helpers"
	"minitwit/src/metrics"
	"minitwit/src/models"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := userRepo.GetByField(ctx, "username", username)
	if err != nil {
		slog.ErrorContext(ctx, "User not found", slog.Any("username", username))
		return nil, err
	}
	return user, nil
}

func GetUserByID(c echo.Context, id int) (*models.User, error) {
	user, err := userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "User not found", slog.Any("user-id", id))
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
		"follower_id":  sessionUserID,
		"following_id": profileUserID,
	}
	followers, err := followerRepo.GetFiltered(c.Request().Context(), conditions, 1, "")

	return err == nil && len(followers) > 0
}

func GetCurrentUser(c echo.Context) (*models.User, error) {
	id, err := helpers.GetSessionUserID(c)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "GetCurrentUser: getSessionUserID returned an error", err)
		return nil, err
	}

	user, err := userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "GetCurrentUser: userRepo.GetByID returned an error", err)
		return nil, err
	}
	return user, nil
}

func CreateUser(username string, email string, hash string) error {
	newUser := helpers.NewUser(username, email, hash)
	err := userRepo.Create(context.Background(), newUser)
	if err != nil {
		utils.LogError("userRepo.Create returned an error", err)
		return err
	}
	metrics.NewUsers.Inc()
	return nil
}

func CountAllUsers(c context.Context) (int, error) {
	user, err := userRepo.CountAll(c)
	if err != nil {
		utils.LogErrorContext(c, "❌ Repository Error: CountAllUsers returned error", err)
		return 0, err
	}
	return user, nil
}

func GetUserMap(c echo.Context, userIDs []int) map[string]*models.User {
	return nil
}
