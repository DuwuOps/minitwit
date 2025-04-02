package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"minitwit/src/handlers/helpers"

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

	user, err := getUserByUsername(c.Request().Context(), username)
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
		log.Printf("User \"/%v\" has requested to follow \"/%v\"\n", username, followsUsername)
		follow, err := getUserByUsername(c.Request().Context(), followsUsername)
		if err != nil {
			log.Printf("getUserIdreturned error: %v\n", err)
			return err
		}

		followerRepo.Create(c.Request().Context(), newFollower(user.UserID, follow.UserID))

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodPost && unfollowsUsername != "" {
		log.Printf("User \"/%v\" has requested to unfollow \"/%v\"\n", username, unfollowsUsername)
		unfollow, err := getUserByUsername(c.Request().Context(), unfollowsUsername)
		if err != nil {
			log.Printf("getUserId returned error: %v\n", err)
			return err
		}

		conditions := map[string]any{
			"who_id": user.UserID, 
			"WHOM_ID": unfollow.UserID,
		}
		followerRepo.DeleteByFields(c.Request().Context(), conditions)

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodGet {
		noFollowersStr := c.QueryParam("no")
		noFollowers := 100
		if noFollowersStr != "" {
			val, err := strconv.Atoi(noFollowersStr)
			if err == nil {
				noFollowers = val
			}
		}

		conditions := map[string]any{
			"who_id": user.UserID,
		}
		followers, err := followerRepo.GetFiltered(c.Request().Context(), conditions, noFollowers, "")
	
		if err != nil {
			log.Printf("Follow: Error retrieving followers for userID=%d: %v", user.UserID, err)
			return err
		}

		var followList []string
		for _, follower := range followers {
			targetUser, err := getUserByID(c.Request().Context(), follower.WhomID)
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

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Printf("FollowUser: getUserByUsername returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}
	
	follower := newFollower(sessionUserId, user.UserID)
	followerRepo.Create(c.Request().Context(), follower)

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

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Printf("row.Scan returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}
	
	conditions := map[string]any{
		"who_id":  sessionUserId,
		"whom_id": user.UserID,
	}
	followerRepo.DeleteByFields(c.Request().Context(), conditions)

	err = helpers.AddFlash(c, fmt.Sprintf("You are no longer following \"%s\"", username))
	if err != nil {
		log.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}
