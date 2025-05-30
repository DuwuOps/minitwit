package handlers

import (
	"net/http"

	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

var PER_PAGE = 30

func AddMessage(c echo.Context) error {
	utils.LogRouteStart(c, "AddMessage", "/add_message")

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		err := c.String(http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			utils.LogErrorEchoContext(c, "echo.Context.String returned an error", err)
		}
	}
	text := c.FormValue("text")
	userId, err := helpers.GetSessionUserID(c)
	if err != nil {
		utils.LogError("getSessionUserID returned an error", err)
		return err
	}

	_ = repo_wrappers.CreateMessage(c, userId, text)

	err = helpers.AddFlash(c, "Your message was recorded")
	if err != nil {
		utils.LogError("addFlash returned an error", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Messages(c echo.Context) error {
	utils.LogRouteStart(c, "Messages", "/msgs")

	err := repo_wrappers.UpdateLatest(c)
	if err != nil {
		utils.LogError("repo_wrappers.UpdateLatest returned an error", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		return err
	}

	noMsgs := GetNumber(c)

	if c.Request().Method == http.MethodGet {
		conditions := map[string]any{
			"flagged": 0,
		}

		msgs, err := repo_wrappers.GetMessagesFiltered(c, conditions, noMsgs)
		if err != nil {
			utils.LogErrorContext(c.Request().Context(), "Messages: repo_wrappers.GetMessagesFiltered returned an error", err)
			return err
		}

		enrichedMsgs := repo_wrappers.EnhanceMessages(c, msgs, true)

		return c.JSON(http.StatusOK, enrichedMsgs)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func MessagesPerUser(c echo.Context) error {
	username := c.Param("username")
	utils.LogRouteStart(c, "MessagesPerUser", "/msgs/:username")

	err := repo_wrappers.UpdateLatest(c)
	if err != nil {
		utils.LogError("repo_wrappers.UpdateLatest returned an error", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		return err
	}

	noMsgs := GetNumber(c)

	user, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {
		conditions := map[string]any{
			"flagged":   0,
			"author_id": user.UserID,
		}

		msgs, err := repo_wrappers.GetMessagesFiltered(c, conditions, noMsgs)
		if err != nil {
			utils.LogErrorContext(c.Request().Context(), "MessagesPerUser: repo_wrappers.GetMessagesFiltered returned an error", err)
			return err
		}

		enrichedMsgs := repo_wrappers.EnhanceMessages(c, msgs, true)

		return c.JSON(http.StatusOK, enrichedMsgs)
	} else if c.Request().Method == http.MethodPost {
		payload, err := helpers.ExtractJson(c)
		if err != nil {
			utils.LogErrorContext(c.Request().Context(), "MessagesPerUser: ExtractJson returned an error", err)
		}

		var requestData string

		if payload != nil {
			requestData = helpers.GetStringValue(payload, "content")
		} else {
			requestData = c.FormValue("content")
		}

		err = repo_wrappers.CreateMessage(c, user.UserID, requestData)
		if err != nil {
			utils.LogErrorContext(c.Request().Context(), "MessagesPerUser: repo_wrappers.CreateMessage returned an error", err)
		}

		return c.JSON(http.StatusNoContent, nil)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func UserTimeline(c echo.Context) error {
	username := c.Param("username")
	utils.LogRouteStart(c, "UserTimeline", "/:username")

	requestedUser, err := repo_wrappers.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		utils.LogError("getUserByUsername returned an error", err)
		err := c.String(http.StatusNotFound, "Not found")
		if err != nil {
			utils.LogErrorEchoContext(c, "echo.Context.String returned an error", err)
		}
	}

	followed := false
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		followed = repo_wrappers.IsFollowingUser(c, requestedUser.UserID)
	}

	conditions := map[string]any{
		"author_id": requestedUser.UserID,
	}

	msgs, err := repo_wrappers.GetMessagesFiltered(c, conditions, PER_PAGE)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "UserTimeline: repo_wrappers.GetMessagesFiltered returned an error", err)
		return err
	}

	enrichedMsgs := repo_wrappers.EnhanceMessages(c, msgs, false)

	user, err := repo_wrappers.GetCurrentUser(c)
	if err != nil {
		utils.LogError("No user found. getCurrentUser returned an error", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		utils.LogError("addFlash returned an error", err)
	}

	data := map[string]any{
		"Messages":    enrichedMsgs,
		"Followed":    followed,
		"ProfileUser": requestedUser,
		"User":        user,
		"Endpoint":    c.Path(),
		"Flashes":     flashes,
	}
	data = utils.MapCSRFToContext(c, data)
	return c.Render(http.StatusOK, "timeline.html", data)
}

func PublicTimeline(c echo.Context) error {
	utils.LogRouteStart(c, "PublicTimeline", "/public")

	conditions := map[string]any{"flagged": 0}
	msgs, err := repo_wrappers.GetMessagesFiltered(c, conditions, PER_PAGE)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "PublicTimeline: repo_wrappers.GetMessagesFiltered returned an error", err)
		return err
	}

	enrichedMsgs := repo_wrappers.EnhanceMessages(c, msgs, false)

	user, err := repo_wrappers.GetCurrentUser(c)
	if err != nil {
		utils.LogError("getCurrentUser returned an error", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		utils.LogError("getFlashes returned an error", err)
	}

	data := map[string]any{
		"Messages": enrichedMsgs,
		"Endpoint": c.Path(),
		"User":     user,
		"Flashes":  flashes,
	}
	data = utils.MapCSRFToContext(c, data)
	return c.Render(http.StatusOK, "timeline.html", data)
}

func Timeline(c echo.Context) error {
	utils.LogRouteStart(c, "Timeline", "/")

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		return c.Redirect(http.StatusFound, "/public")
	}

	sessionUserId, _ := helpers.GetSessionUserID(c)

	conditions := map[string]any{"follower_id": sessionUserId}
	followings, _ := repo_wrappers.GetFollowerFiltered(c, conditions, -1)

	followedUserIDs := []int{sessionUserId}
	for _, f := range followings {
		followedUserIDs = append(followedUserIDs, f.FollowingID)
	}

	conditions = map[string]any{
		"flagged":   0,
		"author_id": followedUserIDs,
	}
	msgs, err := repo_wrappers.GetMessagesFiltered(c, conditions, PER_PAGE)
	if err != nil {
		utils.LogErrorContext(c.Request().Context(), "Timeline: repo_wrappers.GetMessagesFiltered returned an error", err)
		return err
	}

	enrichedMsgs := repo_wrappers.EnhanceMessages(c, msgs, false)

	user, err := repo_wrappers.GetCurrentUser(c)
	if err != nil {
		utils.LogError("No user found. getCurrentUser returned an error", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		utils.LogError("addFlash returned an error", err)
	}

	data := map[string]any{
		"Messages": enrichedMsgs,
		"User":     user,
		"Endpoint": c.Path(),
		"Flashes":  flashes,
	}
	data = utils.MapCSRFToContext(c, data)
	return c.Render(http.StatusOK, "timeline.html", data)
}
