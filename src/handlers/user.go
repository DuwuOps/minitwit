package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

func Follow(c echo.Context) error {
	username := c.Param("username")
	utils.LogRouteStart(c, "Follow", "/fllws/:username")

	err := repo_wrappers.UpdateLatest(c)
	if err != nil {
		utils.LogError("repo_wrappers.UpdateLatest returned an error", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		utils.LogError("notReqFromSimulator returned an error", err)
		return err
	}

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		utils.LogError("getUserId returned an error", err)
		return err
	}

	payload, err := helpers.ExtractJson(c)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "Follow: ExtractJson returned an error", err)
	}
	

	var followsUsername string
	var unfollowsUsername string

	if payload != nil {
		followsUsername = helpers.GetStringValue(payload, "follow")
		unfollowsUsername = helpers.GetStringValue(payload, "unfollow")
	} else {
		followsUsername = c.FormValue("follow")
		unfollowsUsername = c.FormValue("unfollow")
	}

	if c.Request().Method == http.MethodPost && followsUsername != "" {
		slog.InfoContext(c.Request().Context(), "User has requested to follow", slog.Any("username", username), slog.Any("followsUsername", followsUsername))
		follow, err := repo_wrappers.GetUserByUsername(c.Request().Context(), followsUsername)
		if err != nil {
			utils.LogErrorContext(c.Request().Context(), "getUserId returned an error", err)
			return err
		}

		_ = repo_wrappers.CreateFollower(c, user.UserID, follow.UserID)

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodPost && unfollowsUsername != "" {
		slog.InfoContext(c.Request().Context(), "User has requested to unfollow", slog.Any("username", username), slog.Any("followsUsername", unfollowsUsername))
		unfollow, err := repo_wrappers.GetUserByUsername(c.Request().Context(), unfollowsUsername)
		if err != nil {
			utils.LogError("getUserId returned an error", err)
			return err
		}

		_ = repo_wrappers.DeleteFollower(c, user.UserID, unfollow.UserID)

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodGet {
		noFollowers := GetNumber(c)

		conditions := map[string]any{
			"follower_id": user.UserID,
		}
		
		followers, err := repo_wrappers.GetFollowerFiltered(c, conditions, noFollowers)
	
		if err != nil {
			slog.ErrorContext(c.Request().Context(), "Follow: Error retrieving followers for userID", slog.Any("userID", user.UserID), slog.Any("error", err))
			return err
		}

		var followList []string
		for _, follower := range followers {
			targetUser, err := repo_wrappers.GetUserByID(c, follower.FollowingID)
			if err == nil {
				followList = append(followList, targetUser.Username)
			}
		}

		data := map[string]any{
			"follows": followList,
		}
		slog.InfoContext(c.Request().Context(), "", slog.Any("data", data))

		return c.JSON(http.StatusOK, data)
	}

	slog.Error("\"/fllws/:username\" was entered wrongly!")
	return c.JSON(http.StatusBadRequest, nil)
}

func FollowUser(c echo.Context) error {
	username := c.Param("username")
	utils.LogRouteStart(c, "FollowUser", "/:username/follow")

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		err := c.String(http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			utils.LogErrorEchoContext(c, "echo.Context.String returned an error", err)
		}
	}

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "FollowUser: getUserByUsername returned an error", err)
		err := c.String(http.StatusNotFound, "Not found")
		if err != nil {
			utils.LogErrorEchoContext(c, "echo.Context.String returned an error", err)
		}
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		utils.LogError("getSessionUserID returned an error", err)
		return err
	}
	
	_ = repo_wrappers.CreateFollower(c, sessionUserId, user.UserID)

	err = helpers.AddFlash(c, fmt.Sprintf("You are now following \"%s\"", username))
	if err != nil {
		utils.LogError("addFlash returned an error", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}

func UnfollowUser(c echo.Context) error {
	username := c.Param("username")
	utils.LogRouteStart(c, "UnfollowUser", "/:username/unfollow")

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		err := c.String(http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			utils.LogErrorEchoContext(c, "echo.Context.String returned an error", err)
		}
	}

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		utils.LogError("row.Scan returned an error", err)
		err := c.String(http.StatusNotFound, "Not found")
		if err != nil {
			utils.LogErrorEchoContext(c, "echo.Context.String returned an error", err)
		}
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		utils.LogError("getSessionUserID returned an error", err)
		return err
	}

	_ = repo_wrappers.DeleteFollower(c, sessionUserId, user.UserID)

	err = helpers.AddFlash(c, fmt.Sprintf("You are no longer following \"%s\"", username))
	if err != nil {
		utils.LogError("addFlash returned an error", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}
