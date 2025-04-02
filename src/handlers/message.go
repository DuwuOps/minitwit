package handlers

import (
	"log"
	"net/http"
	"strconv"

	"minitwit/src/handlers/helpers"

	"github.com/labstack/echo/v4"
)

var PER_PAGE = 30

func AddMessage(c echo.Context) error {
	log.Println("🎺 User entered AddMessage via route \"/add_message\"")

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}
	text := c.FormValue("text")
	userId, err := helpers.GetSessionUserID(c)
	if err != nil {
		log.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}

	newMessage := newMessage(userId, text)
	messageRepo.Create(c.Request().Context(), newMessage)

	err = helpers.AddFlash(c, "Your message was recorded")
	if err != nil {
		log.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Messages(c echo.Context) error {
	log.Println("🎺 User entered Messages via route \"/msgs\"")

	err := helpers.UpdateLatest(c)
	if err != nil {
		log.Printf("helpers.UpdateLatest returned error: %v\n", err)
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
			log.Printf("messages: messageRepo.GetFiltered returned error: %v\n", err)
			return err
		}
	
		enrichedMsgs := []map[string]any{}
		for _, msg := range msgs {
			enrichedMsg := map[string]any{
				"pub_date": msg.PubDate,
				"content": msg.Text,
			}
			
			author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
			if author != nil {
				enrichedMsg["user"] = author.Username
				enrichedMsg["email"] = author.Email
			} else {
				log.Printf("⚠️ Warning: Could not find user for message author_id=%d\n", msg.AuthorID)
				enrichedMsg["user"] = "Unknown"
				enrichedMsg["email"] = ""
			}
	
			enrichedMsgs = append(enrichedMsgs, enrichedMsg)
		}

		return c.JSON(http.StatusOK, enrichedMsgs)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func MessagesPerUser(c echo.Context) error {
	username := c.Param("username")
	log.Printf("🎺 User entered MessagesPerUser via route \"/msgs/:username\" as \"/%v\" and HTTP method %v\n", username, c.Request().Method)

	err := helpers.UpdateLatest(c)
	if err != nil {
		log.Printf("helpers.UpdateLatest returned error: %v\n", err)
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
			log.Printf("messages: messageRepo.GetFiltered returned error: %v\n", err)
			return err
		}

		enrichedMsgs := []map[string]any{}
		for _, msg := range msgs {
			enrichedMsg := map[string]any{
				"pub_date": msg.PubDate,
				"content": msg.Text,
			}
			
			author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
			if author != nil {
				enrichedMsg["user"] = author.Username
			} else {
				log.Printf("⚠️ Warning: Could not find user for message author_id=%d\n", msg.AuthorID)
				enrichedMsg["user"] = "Unknown"
			}
	
			enrichedMsgs = append(enrichedMsgs, enrichedMsg)
		}

		return c.JSON(http.StatusOK, enrichedMsgs)
	} else if c.Request().Method == http.MethodPost {
		payload, err := helpers.ExtractJson(c)
		if err != nil {
			log.Printf("MessagesPerUser: ExtractJson returned error: %v\n", err)
		}

		var requestData string

		if payload != nil {
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
	log.Printf("🎺 User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	requestedUser, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Printf("getUserByUsername returned error: %v\n", err)
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
		log.Printf("UserTimeline: messageRepo.GetFiltered returned error: %v\n", err)
		return err
	}

	enrichedMsgs := []map[string]any{}
	for _, msg := range msgs {
		enrichedMsg := map[string]any{
			"pub_date": msg.PubDate,
			"text": msg.Text,
		}
		
		author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
		if author != nil {
			enrichedMsg["username"] = author.Username
			enrichedMsg["email"] = author.Email
		} else {
			log.Printf("⚠️ Warning: Could not find user for message author_id=%d\n", msg.AuthorID)
			enrichedMsg["username"] = "Unknown"
			enrichedMsg["email"] = ""
		}

		enrichedMsgs = append(enrichedMsgs, enrichedMsg)
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Printf("No user found. getCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		log.Printf("addFlash returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages":    enrichedMsgs,
		"Followed":    followed,
		"ProfileUser": requestedUser,
		"User":        user,
		"Endpoint":    c.Path(),
		"Flashes":     flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}

func PublicTimeline(c echo.Context) error {
	log.Println("🎺 User entered PublicTimeline via route \"/public\"")
	
	conditions := map[string]any{"flagged": 0}
	msgs, err := messageRepo.GetFiltered(c.Request().Context(), conditions, PER_PAGE, "pub_date DESC")
	if err != nil {
		log.Printf("PublicTimeline: messageRepo.GetFiltered returned error: %v\n", err)
	}

	enrichedMsgs := []map[string]any{}
	for _, msg := range msgs {
		enrichedMsg := map[string]any{
			"pub_date": msg.PubDate,
			"text": msg.Text,
		}

		author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
		if author != nil {
			enrichedMsg["username"] = author.Username
			enrichedMsg["email"] = author.Email
		} else {
			log.Printf("⚠️ Warning: Could not find user for message author_id=%d\n", msg.AuthorID)
			enrichedMsg["username"] = "Unknown"
			enrichedMsg["email"] = ""
		}

		enrichedMsgs = append(enrichedMsgs, enrichedMsg)
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Printf("getCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		log.Printf("getFlashes returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages": enrichedMsgs,
		"Endpoint": c.Path(),
		"User":     user,
		"Flashes":  flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}

func Timeline(c echo.Context) error {
	log.Println("🎺 User entered Timeline via route \"/\"")

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		return c.Redirect(http.StatusFound, "/public")
	}

	sessionUserId, _ := helpers.GetSessionUserID(c)

	conditions := map[string]any{"who_id": sessionUserId}
	followings, _ := followerRepo.GetFiltered(c.Request().Context(), conditions, -1, "")

	followedUserIDs := []int{sessionUserId} 
	for _, f := range followings {
		followedUserIDs = append(followedUserIDs, f.WhomID)
	}

	conditions = map[string]any{
		"flagged": 0,
		"author_id": followedUserIDs,
	}
	msgs, err := messageRepo.GetFiltered(c.Request().Context(), conditions, PER_PAGE, "pub_date DESC")
	if err != nil {
		log.Printf("Timeline: messageRepo.GetFiltered returned error: %v\n", err)
		return err
	}

	var enrichedMsgs []map[string]any
    for _, msg := range msgs {
		enrichedMsg := map[string]any{
			"pub_date": msg.PubDate,
			"text": msg.Text,
		}

		author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
		if author != nil {
			enrichedMsg["username"] = author.Username
			enrichedMsg["email"] = author.Email
		} else {
			log.Printf("⚠️ Warning: Could not find user for message author_id=%d\n", msg.AuthorID)
			enrichedMsg["username"] = "Unknown"
			enrichedMsg["email"] = ""
		}

		enrichedMsgs = append(enrichedMsgs, enrichedMsg)
    }

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Printf("No user found. getCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		log.Printf("addFlash returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages": enrichedMsgs,
		"User":     user,
		"Endpoint": c.Path(),
		"Flashes":  flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}
