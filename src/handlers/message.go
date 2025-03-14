package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"context"

	"minitwit/src/datalayer"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"

	"github.com/labstack/echo/v4"
)

var messageRepo *datalayer.Repository[models.Message]

func SetMessageRepo(repo *datalayer.Repository[models.Message]) {
	messageRepo = repo
}

var followerRepo *datalayer.Repository[models.Follower]

func SetFollowerRepo(repo *datalayer.Repository[models.Follower]) {
	followerRepo = repo
}

var PER_PAGE = 30

func AddMessage(c echo.Context) error {
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	text := c.FormValue("text")
	userId, err := helpers.GetSessionUserID(c)
	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}

	newMessage := newMessage(userId, text)

	err = messageRepo.Create(c.Request().Context(), newMessage)

	err = helpers.AddFlash(c, "Your message was recorded")
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Messages(c echo.Context) error {
    log.Println("User entered Messages via route \"/:msgs\"")

    if err := helpers.ValidateRequest(c); err != nil {
        log.Printf("Request blocked by NotReqFromSimulator: %v", err)
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden"})
    }

    messages, err := handleGetMessages(c, nil, true)
    if err != nil {
        log.Printf("Error fetching messages: %v", err)
        return err
    }

    return c.JSON(http.StatusOK, messages)
}

func MessagesPerUser(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered MessagesPerUser via route \"/msgs/:username\" as \"/%v\"\n", username)

	if err := helpers.ValidateRequest(c); err != nil {
		return err
	}

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	if c.Request().Method == http.MethodGet {
		messages, err := handleGetMessages(c, user, true)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, messages)
	} else if c.Request().Method == http.MethodPost {
		return handlePostMessage(c, user)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func UserTimeline(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	user, err := getUserByUsername(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	followed := isFollowingUser(c, user.UserID)

	messages, err := handleGetMessages(c, user, false)
	if err != nil {
		return err
	}

	log.Printf("📝 Messages for user %s: %v", user.Username, messages)

	loggedInUser, err := GetCurrentUser(c)
	if err != nil {
		log.Printf("No user found. GetCurrentUser returned error: %v\n", err)
	}

	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	data := map[string]any{
		"Messages":    messages,
		"Followed":    followed,
		"ProfileUser": user,
		"User":        loggedInUser,
		"Endpoint":    c.Path(),
		"Flashes":     flashes,
	}

	log.Printf("📤 Data passed to template: %v", data)

	return c.Render(http.StatusOK, "timeline.html", data)
}


func PublicTimeline(c echo.Context) error {
	log.Println("User entered PublicTimeline via route \"/public\"")

	conditions := map[string]any{"flagged": 0} 
	return handleRenderTimeline(c, conditions, nil)
}

func Timeline(c echo.Context) error {
    log.Println("User entered Timeline via route \"/\"")
    log.Printf("We got a visitor from: %s", c.Request().RemoteAddr)

    if loggedIn, _ := helpers.IsUserLoggedIn(c); !loggedIn {
        return c.Redirect(http.StatusFound, "/public")
    }

    sessionUserId, err := helpers.GetSessionUserID(c)
    if err != nil {
        log.Printf("Failed to get session user ID: %v\n", err)
        return err
    }

    // Get list of user IDs the logged-in user follows
	conditions := map[string]any{"who_id": sessionUserId}
	followers, err := followerRepo.GetFiltered(c.Request().Context(), conditions, -1, "")
	if err != nil {
		log.Printf("Error fetching followers: %v\n", err)
		return err
	}

	// Only include followed users + logged-in user in timeline
	followedUserIDs := []int{sessionUserId}  // Always include self
	for _, f := range followers {
		followedUserIDs = append(followedUserIDs, f.WhomID)
	}

	log.Printf("🔍 Fetching messages for users: %v", followedUserIDs)

	return handleRenderTimeline(c, map[string]any{
		"flagged": 0,
		"author_id": followedUserIDs, 
	}, nil)	
}


func handleRenderTimeline(c echo.Context, conditions map[string]any, user *models.User) error {
	messages, err := messageRepo.GetFiltered(c.Request().Context(), conditions, PER_PAGE, "pub_date DESC")
	if err != nil {
		log.Printf("❌ Error retrieving messages: %v", err)
		return err
	}

	var enrichedMessages []map[string]any
	for _, msg := range messages {
		author, err := getUserByID(c.Request().Context(), msg.AuthorID)
		if err != nil {
			log.Printf("⚠️ Warning: Could not find user for message author_id=%d", msg.AuthorID)
			continue
		}
		
		enrichedMessages = append(enrichedMessages, map[string]any{
			"text":     msg.Text,
			"pub_date": msg.PubDate,
			"username": author.Username,
		})
		

		/*enrichedMessages = append(enrichedMessages, map[string]any{
			"user":     author.Username,
			"content":  msg.Text,
			"pub_date": msg.PubDate,
		})*/
	}

	log.Printf("📥 Filtered Messages Before Rendering: %+v", enrichedMessages)

	user, _ = GetCurrentUser(c)
	flashes, _ := getFlashes(c)

	data := map[string]any{
		"Messages": enrichedMessages,
		"Endpoint": c.Path(),
		"User":     user,
		"Flashes":  flashes,
	}

	return c.Render(http.StatusOK, "timeline.html", data)
}


func isFollowingUser(c echo.Context, profileUserID int) bool {
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

func newMessage(authorID int, text string) *models.Message {
	return &models.Message{
		AuthorID: authorID,
		Text:     text,
		PubDate:  time.Now().Unix(),
		Flagged:  0,
	}
}

func getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := userRepo.GetByField(ctx, "username", username)
	if err != nil {
		log.Printf("User not found: %s", username)
		return nil, err
	}
	return user, nil
}

func parseMessageLimit(c echo.Context) int {
	noMsgs := 100
	if noMsgsStr := c.QueryParam("no"); noMsgsStr != "" {
		if val, err := strconv.Atoi(noMsgsStr); err == nil {
			noMsgs = val
		}
	}
	return noMsgs
}

func handleGetMessages(c echo.Context, user *models.User, useContentKey bool) ([]map[string]any, error) {
	noMsgs := parseMessageLimit(c)

	conditions := map[string]any{
		"flagged": 0,
	}

	if user != nil {
		conditions["author_id"] = user.UserID
		log.Printf("🔍 Fetching messages for user: %s (ID: %d)", user.Username, user.UserID)
	} else {
		log.Println("🔍 Fetching all public messages")
	}

	messages, err := messageRepo.GetFiltered(c.Request().Context(), conditions, noMsgs, "pub_date DESC")
	if err != nil {
		log.Printf("❌ Error retrieving messages: %v", err)
		return nil, err
	}

	var filteredMsgs []map[string]any
	for _, msg := range messages {
		username := "Unknown"

		author, _ := userRepo.GetByID(c.Request().Context(), msg.AuthorID)
		if author != nil {
			username = author.Username
		}

		messageData := map[string]any{
			"pub_date": msg.PubDate,
		}

		if useContentKey {
			messageData["content"] = msg.Text
			messageData["user"] = username
		} else {
			messageData["text"] = msg.Text
			messageData["username"] = username
		}

		filteredMsgs = append(filteredMsgs, messageData)
	}

	log.Printf("📥 Filtered Messages: %v", filteredMsgs)
	return filteredMsgs, nil
}



func getFlashes(c echo.Context) ([]string, error) {
	flashes, err := helpers.GetFlashes(c)
	if err != nil {
		log.Printf("getFlashes returned error: %v", err)
		return nil, err
	}
	return flashes, nil
}

func getUserByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("User not found for ID: %d", userID)
		return nil, err
	}
	return user, nil
}

func handlePostMessage(c echo.Context, user *models.User) error {
	requestData, err := extractMessageContent(c)
	if err != nil {
		log.Printf("Error extracting message content: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid message content"})
	}

	newMessage := newMessage(user.UserID, requestData)
	err = messageRepo.Create(c.Request().Context(), newMessage)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		return err
	}

	return c.JSON(http.StatusNoContent, nil)
}

func extractMessageContent(c echo.Context) (string, error) {
	payload, err := helpers.ExtractJson(c)
	if err == nil {
		return helpers.GetStringValue(payload, "content"), nil
	}
	return c.FormValue("content"), nil
}

