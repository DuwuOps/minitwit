package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"minitwit/src/datalayer"
	"minitwit/src/handlers/helpers"
	"minitwit/src/models"

	"github.com/labstack/echo/v4"
)

var PER_PAGE = 30

func AddMessage(c echo.Context, db *sql.DB) error {
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

	db.Exec(`insert into message (author_id, text, pub_date, flagged)
			 values (?, ?, ?, 0)`,
		userId, (text + "okay let's go"), time.Now().Unix(),
	)

	err = helpers.AddFlash(c, "Your message was recorded")
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Messages(c echo.Context, db *sql.DB) error {
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
		rows, err := datalayer.QueryDB(db, `SELECT message.*, user.* FROM message, user
					WHERE message.flagged = 0 AND message.author_id = user.user_id
					ORDER BY message.pub_date DESC LIMIT ?`,
			noMsgs,
		)
		if err != nil {
			fmt.Printf("messages: queryDB returned error: %v\n", err)
			return err
		}

		msgs, err := helpers.RowsToMapList(rows)
		if err != nil {
			fmt.Printf("messages: rowsToMapList returned error: %v\n", err)
			return err
		}

		filteredMsgs := []map[string]any{}
		for _, msg := range msgs {
			filteredMsg := map[string]any{
				"content":  msg["text"],
				"pub_date": msg["pub_date"],
				"user":     msg["username"],
			}
			filteredMsgs = append(filteredMsgs, filteredMsg)
		}

		return c.JSON(http.StatusOK, filteredMsgs)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func MessagesPerUser(c echo.Context, db *sql.DB) error {
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

	userId, err := datalayer.GetUserId(username, db)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {

		rows, err := datalayer.QueryDB(db, `SELECT message.*, user.* FROM message, user
					WHERE message.flagged = 0 AND
					user.user_id = message.author_id AND user.user_id = ?
					ORDER BY message.pub_date DESC LIMIT ?`,
			userId, noMsgs,
		)
		if err != nil {
			fmt.Printf("messages: queryDB returned error: %v\n", err)
			return err
		}

		msgs, err := helpers.RowsToMapList(rows)
		if err != nil {
			fmt.Printf("messages: rowsToMapList returned error: %v\n", err)
			return err
		}

		filteredMsgs := []map[string]any{}
		for _, msg := range msgs {
			filteredMsg := map[string]any{
				"content":  msg["text"],
				"pub_date": msg["pub_date"],
				"user":     msg["username"],
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

		fmt.Printf("requestData: %v\n", requestData)
		query := `INSERT INTO message (author_id, text, pub_date, flagged)
                   VALUES (?, ?, ?, 0)`

		db.Exec(query,
			userId, requestData, noMsgs, time.Now().Unix(),
		)

		return c.JSON(http.StatusNoContent, nil)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func UserTimeline(c echo.Context, db *sql.DB) error {
	username := c.Param("username")
	fmt.Printf("User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	row := datalayer.QueryDbSingle(db, "select * from user where username = ?", username)
	var requestedUser models.User
	err := row.Scan(&requestedUser.UserID, &requestedUser.Username, &requestedUser.Email, &requestedUser.PwHash)
	if err != nil {
		fmt.Printf("row.Scan returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	followed := false
	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if loggedIn {
		sessionUserId, _ := helpers.GetSessionUserID(c)
		follow_result := datalayer.QueryDbSingle(db, `select 1 from follower where
             follower.who_id = ? and follower.whom_id = ?`,
			sessionUserId, requestedUser.UserID)

		// The query should return a 1, if the user follows the user of the timeline.
		var result int
		err := follow_result.Scan(&result)
		followed = err == nil
	}

	rows, err := datalayer.QueryDB(db, `select message.*, user.* from message, user where
                            user.user_id = message.author_id and user.user_id = ?
                            order by message.pub_date desc limit ?`,
		requestedUser.UserID, PER_PAGE,
	)

	if err != nil {
		fmt.Printf("UserTimeline: queryDB returned error: %v\n", err)
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
		"Messages":    msgs,
		"Followed":    followed,
		"ProfileUser": requestedUser,
		"User":        user,
		"Endpoint":    c.Path(),
		"Flashes":     flashes,
	}
	return c.Render(http.StatusOK, "timeline.html", data)
}

func PublicTimeline(c echo.Context, db *sql.DB) error {
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

func Timeline(c echo.Context, db *sql.DB) error {
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
