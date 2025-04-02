package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"minitwit/src/handlers/helpers"

	"github.com/labstack/echo/v4"
)

var PER_PAGE = 30

func AddMessage(c echo.Context) error {
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}
	text := c.FormValue("text")
	userId, err := helpers.GetSessionUserID(c)
	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}

	newMessage := newMessage(userId, text)
	messageRepo.Create(c.Request().Context(), newMessage)

	err = helpers.AddFlash(c, "Your message was recorded")
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Messages(c echo.Context) error {
	fmt.Printf("User entered Messages via route \"/:msgs\"")

	err := helpers.UpdateLatest(c)
	if err != nil {
		fmt.Printf("helpers.UpdateLatest returned error: %v\n", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		return err
	}

	noMsgsStr := c.QueryParam("no")
	noMsgs := 100
	if noMsgsStr != "" {
		val, err := strconv.Atoi(noMsgsStr)
		if err == nil {
			noMsgs = val
		}
	}

	if c.Request().Method == http.MethodGet {
		conditions := map[string]any{
			"flagged": 0,
		}
		msgs, err := messageRepo.GetFiltered(c.Request().Context(), conditions, noMsgs, "pub_date DESC")
		if err != nil {
			fmt.Printf("messages: messageRepo.GetFiltered returned error: %v\n", err)
			return err
		}

		filteredMsgs := []map[string]any{}
		for _, msg := range msgs {
			filteredMsg := map[string]any{
				"pub_date": msg.PubDate,
				"content": msg.Text,
			}
			
			author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
			if author != nil {
				filteredMsg["user"] = author.Username
			} else {
				filteredMsg["user"] = "Unknown"
			}
	
			filteredMsgs = append(filteredMsgs, filteredMsg)
		}

		return c.JSON(http.StatusOK, filteredMsgs)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func MessagesPerUser(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered MessagesPerUser via route \"/msgs/:username\" as \"/%v\"\n", username)

	err := helpers.UpdateLatest(c)
	if err != nil {
		fmt.Printf("helpers.UpdateLatest returned error: %v\n", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		return err
	}

	noMsgsStr := c.QueryParam("no")
	noMsgs := 100
	if noMsgsStr != "" {
		val, err := strconv.Atoi(noMsgsStr)
		if err == nil {
			noMsgs = val
		}
	}

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {
		conditions := map[string]any{
			"flagged": 0,
			"author_id": user.UserID,
		}

		msgs, err := messageRepo.GetFiltered(c.Request().Context(), conditions, noMsgs, "pub_date DESC")
		if err != nil {
			fmt.Printf("messages: messageRepo.GetFiltered returned error: %v\n", err)
			return err
		}

		filteredMsgs := []map[string]any{}
		for _, msg := range msgs {
			filteredMsg := map[string]any{
				"pub_date": msg.PubDate,
				"content": msg.Text,
			}
			
			author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
			if author != nil {
				filteredMsg["user"] = author.Username
			} else {
				filteredMsg["user"] = "Unknown"
			}
	
			filteredMsgs = append(filteredMsgs, filteredMsg)
		}

		return c.JSON(http.StatusOK, filteredMsgs)
	} else if c.Request().Method == http.MethodPost {
		payload, err := helpers.ExtractJson(c)

		var requestData string

		if err == nil {
			requestData = helpers.GetStringValue(payload, "content")
		} else {
			requestData = c.FormValue("content")
		}

		newMessage := newMessage(user.UserID, requestData)
		messageRepo.Create(c.Request().Context(), newMessage)

		return c.JSON(http.StatusNoContent, nil)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func UserTimeline(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	requestedUser, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		fmt.Printf("getUserByUsername returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	followed := false
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		followed = isFollowingUser(c, requestedUser.UserID)
	}


	conditions := map[string]any{
		"author_id": requestedUser.UserID,
	}

	msgs, err := messageRepo.GetFiltered(c.Request().Context(), conditions, PER_PAGE, "pub_date DESC")
	if err != nil {
		fmt.Printf("UserTimeline: messageRepo.GetFiltered returned error: %v\n", err)
		return err
	}

	filteredMsgs := []map[string]any{}
	for _, msg := range msgs {
		filteredMsg := map[string]any{
			"pub_date": msg.PubDate,
			"text": msg.Text,
		}
		
		author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
		if author != nil {
			filteredMsg["username"] = author.Username
		} else {
			filteredMsg["username"] = "Unknown"
		}

		filteredMsgs = append(filteredMsgs, filteredMsg)
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		fmt.Printf("No user found. getCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages":    filteredMsgs,
		"Followed":    followed,
		"ProfileUser": requestedUser,
		"User":        user,
		"Endpoint":    c.Path(),
		"Flashes":     flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}

func PublicTimeline(c echo.Context) error {
	log.Println("User entered PublicTimeline via route \"/public\"")

	rows, err := datalayer.QueryDB(db, `select message.*, user.* from message, user
                            where message.flagged = 0 and message.author_id = user.user_id
                            order by message.pub_date desc limit ?`,
		PER_PAGE,
	)
	if err != nil {
		fmt.Printf("PublicTimeline: queryDB returned error: %v\n", err)
		return err
	}

	msgs, err := helpers.RowsToMapList(rows)
	if err != nil {
		fmt.Printf("rowsToMapList returned error: %v\n", err)
		return err
	}

	user, err := GetCurrentUser(c, db)
	if err != nil {
		fmt.Printf("getCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		fmt.Printf("getFlashes returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages": msgs,
		"Endpoint": c.Path(),
		"User":     user,
		"Flashes":  flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}

func Timeline(c echo.Context) error {
	log.Println("User entered Timeline via route \"/\"")
	log.Printf("We got a visitor from: %s", c.Request().RemoteAddr)
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		return c.Redirect(http.StatusFound, "/public")
	}

	sessionUserId, _ := helpers.GetSessionUserID(c)
	rows, err := datalayer.QueryDB(db, `select message.*, user.* from message, user
                          where message.flagged = 0 and message.author_id = user.user_id and (
                              user.user_id = ? or
                              user.user_id in (select whom_id from follower
                                                      where who_id = ?))
                          order by message.pub_date desc limit ?`,
		sessionUserId, sessionUserId, PER_PAGE,
	)

	if err != nil {
		fmt.Printf("Timeline: queryDB returned error: %v\n", err)
		return err
	}

	msgs, err := helpers.RowsToMapList(rows)
	if err != nil {
		fmt.Printf("rowsToMapList returned error: %v\n", err)
		return err
	}

	user, err := GetCurrentUser(c, db)
	if err != nil {
		fmt.Printf("No user found. getCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages": msgs,
		"User":     user,
		"Endpoint": c.Path(),
		"Flashes":  flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}
