package handlers

import (
	"fmt"
	"log"
	"net/http"

	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"

	"github.com/labstack/echo/v4"
)

func Follow(c echo.Context) error {
	username := c.Param("username")
	log.Printf("🎺 User entered Follow via route \"/fllws/:username\" as \"/%v\"\n", username)

	err := helpers.UpdateLatest(c)
	if err != nil {
		log.Printf("helpers.UpdateLatest returned error: %v\n", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		log.Printf("notReqFromSimulator returned error: %v\n", err)
		return err
	}

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Printf("getUserId returned error: %v\n", err)
		return err
	}

	payload, err := helpers.ExtractJson(c)
	if err != nil {
		log.Printf("Follow: ExtractJson returned error: %v\n", err)
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
		log.Printf("User \"%v\" has requested to follow \"%v\"\n", username, followsUsername)
		follow, err := repo_wrappers.GetUserByUsername(c.Request().Context(), followsUsername)
		if err != nil {
			log.Printf("getUserIdreturned error: %v\n", err)
			return err
		}

		_ = repo_wrappers.CreateFollower(c, user.UserID, follow.UserID)

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodPost && unfollowsUsername != "" {
		log.Printf("User \"%v\" has requested to unfollow \"%v\"\n", username, unfollowsUsername)
		unfollow, err := repo_wrappers.GetUserByUsername(c.Request().Context(), unfollowsUsername)
		if err != nil {
			log.Printf("getUserId returned error: %v\n", err)
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
			log.Printf("Follow: Error retrieving followers for userID=%d: %v", user.UserID, err)
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
		log.Printf("data: %v\n", data)

		return c.JSON(http.StatusOK, data)
	}

	log.Printf("ERROR: \"/fllws/:username\" was entered wrongly!\n")
	return c.JSON(http.StatusBadRequest, nil)
}

func FollowUser(c echo.Context) error {
	username := c.Param("username")
	log.Printf("🎺 User entered FollowUser via route \"/:username/follow\" as \"/%v\"\n", username)

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Printf("FollowUser: getUserByUsername returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}
	
	_ = repo_wrappers.CreateFollower(c, sessionUserId, user.UserID)

	err = helpers.AddFlash(c, fmt.Sprintf("You are now following \"%s\"", username))
	if err != nil {
		log.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}

func UnfollowUser(c echo.Context) error {
	username := c.Param("username")
	log.Printf("🎺 User entered UnfollowUser via route \"/:username/unfollow\" as \"/%v\"\n", username)

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Printf("row.Scan returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}

	_ = repo_wrappers.DeleteFollower(c, sessionUserId, user.UserID)

	err = helpers.AddFlash(c, fmt.Sprintf("You are no longer following \"%s\"", username))
	if err != nil {
		log.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}
