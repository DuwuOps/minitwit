package handlers

import (
	"fmt"
	"net/http"
	"log"
	"context"
	"strconv"

	"minitwit/src/handlers/helpers"
	"minitwit/src/models"

	"github.com/labstack/echo/v4"
)

func GetCurrentUser(c echo.Context) (*models.User, error) {
	id, err := helpers.GetSessionUserID(c)
	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return nil, err
	}

	user, err := userRepo.GetByID(c.Request().Context(), id)

	if err != nil {
		log.Printf("User not found in database for userID %d: %v", id, err)
		return nil, err
	}

	fmt.Printf("Found user in database! %v\n", user)
	fmt.Printf("user.UserID: %v\n", user.UserID)
	fmt.Printf("user.Username: %v\n", user.Username)
	fmt.Printf("user.Email: %v\n", user.Email)
	fmt.Printf("user.PwHash: %v\n", user.PwHash)
	return user, nil
}

func Follow(c echo.Context) error {
	username := c.Param("username")
	log.Printf("User entered Follow via route \"/fllws/%s\"", username)

	if err := helpers.ValidateRequest(c); err != nil {
		return err
	}

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	switch c.Request().Method {
	case http.MethodPost:
		return handleFollowAction(c, user.UserID)
	case http.MethodGet:
		return handleGetFollowers(c, user.UserID)
	default:
		log.Printf("Invalid request method for Follow route: %s", c.Request().Method)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request method"})
	}
}

func FollowUser(c echo.Context) error {
	return handleUserFollowAction(c, true)
}

func UnfollowUser(c echo.Context) error {
	return handleUserFollowAction(c, false)
}

func handleUserFollowAction(c echo.Context, follow bool) error {
	username := c.Param("username")
	log.Printf("User entered %s via route \"/%s\"", 
		map[bool]string{true: "FollowUser", false: "UnfollowUser"}[follow], username)

	if err := enforceLogin(c); err != nil {
		return err
	}

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		return err
	}

	if follow {
		err = followerRepo.Create(c.Request().Context(), &models.Follower{
			WhoID:  sessionUserId,
			WhomID: user.UserID,
		})
	} else {
		err = followerRepo.DeleteByFields(c.Request().Context(), map[string]any{
			"who_id":  sessionUserId,
			"whom_id": user.UserID,
		})
	}

	if err != nil {
		log.Printf("Error processing follow/unfollow action for user %s: %v", username, err)
		return err
	}

	message := fmt.Sprintf("You are now %s \"%s\"", 
		map[bool]string{true: "following", false: "no longer following"}[follow], username)
	helpers.AddFlash(c, message)

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}

func enforceLogin(c echo.Context) error {
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	return nil
}

func handleFollowAction(c echo.Context, userID int) error {
	ctx := c.Request().Context()
	followUsername, unfollowUsername := extractFollowRequest(c)

	if followUsername != "" {
		return processFollowAction(ctx, userID, followUsername, true)
	}

	if unfollowUsername != "" {
		return processFollowAction(ctx, userID, unfollowUsername, false)
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid follow request"})
}

func processFollowAction(ctx context.Context, userID int, targetUsername string, follow bool) error {
	targetUser, err := getUserByUsername(ctx, targetUsername)
	if err != nil {
		log.Printf("Follow/unfollow target not found: %s", targetUsername)
		return fmt.Errorf("user not found")
	}

	if follow {
		err = followerRepo.Create(ctx, &models.Follower{
			WhoID:  userID,
			WhomID: targetUser.UserID,
		})
	} else {
		err = followerRepo.DeleteByFields(ctx, map[string]any{
			"who_id":  userID,
			"whom_id": targetUser.UserID,
		})
	}

	if err != nil {
		log.Printf("Error processing follow/unfollow action for %s: %v", targetUsername, err)
		return err
	}

	action := map[bool]string{true: "followed", false: "unfollowed"}[follow]
	log.Printf("User %d %s %s", userID, action, targetUsername)
	return nil 
}


func handleGetFollowers(c echo.Context, userID int) error {
	noFollowers := parseQueryParam(c, "no", 100)

	followers, err := followerRepo.GetFiltered(c.Request().Context(), map[string]any{
		"who_id": userID,
	}, noFollowers, "") 

	if err != nil {
		log.Printf("Error retrieving followers for userID=%d: %v", userID, err)
		return err
	}

	var followerUsernames []string
	for _, follower := range followers {
		targetUser, err := getUserByID(c.Request().Context(), follower.WhomID)
		if err == nil {
			followerUsernames = append(followerUsernames, targetUser.Username)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"follows": followerUsernames})
}


func extractFollowRequest(c echo.Context) (string, string) {
	payload, err := helpers.ExtractJson(c)
	if err == nil {
		return helpers.GetStringValue(payload, "follow"), helpers.GetStringValue(payload, "unfollow")
	}
	return c.FormValue("follow"), c.FormValue("unfollow")
}

func parseQueryParam(c echo.Context, param string, defaultValue int) int {
	valueStr := c.QueryParam(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid query param %s: %v", param, err)
		return defaultValue
	}

	return value
}