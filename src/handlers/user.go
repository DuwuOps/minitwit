package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"minitwit/src/handlers/helpers"

	"github.com/labstack/echo/v4"
)

func GetCurrentUser(c echo.Context, db *sql.DB) (*models.User, error) {
	id, err := helpers.GetSessionUserID(c)
	var user models.User

	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return nil, err
	}

	rows := datalayer.QueryDbSingle(db, "select * from user where user_id = ?",
		id,
	)

	err = rows.Scan(&user.UserID, &user.Username, &user.Email, &user.PwHash)
	if err != nil {
		fmt.Printf("rows.Scan returned error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Found user in database! %v\n", user)
	fmt.Printf("user.UserID: %v\n", user.UserID)
	fmt.Printf("user.Username: %v\n", user.Username)
	fmt.Printf("user.Email: %v\n", user.Email)
	fmt.Printf("user.PwHash: %v\n", user.PwHash)
	return &user, nil
}

func Follow(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered Follow via route \"/fllws/:username\" as \"/%v\"\n", username)

	err := helpers.UpdateLatest(c)
	if err != nil {
		fmt.Printf("helpers.UpdateLatest returned error: %v\n", err)
		return err
	}

	err = helpers.NotReqFromSimulator(c)
	if err != nil {
		fmt.Printf("notReqFromSimulator returned error: %v\n", err)
		return err
	}

	userId, err := datalayer.GetUserId(username, db)
	if err != nil {
		fmt.Printf("getUserId returned error: %v\n", err)
		return err
	}

	payload, err := helpers.ExtractJson(c)

	var followsUsername string
	var unfollowsUsername string

	if err == nil {
		followsUsername = helpers.GetStringValue(payload, "follow")
		unfollowsUsername = helpers.GetStringValue(payload, "unfollow")
	} else {
		followsUsername = c.FormValue("follow")
		unfollowsUsername = c.FormValue("unfollow")
	}

	if c.Request().Method == http.MethodPost && followsUsername != "" {
		fmt.Printf("\"/fllws/:username\" running as a Post-Method, where follow in c.FormParams()")

		followsUserId, err := datalayer.GetUserId(followsUsername, db)
		if err != nil {
			fmt.Printf("getUserIdreturned error: %v\n", err)
			return err
		}

		query := `INSERT INTO follower (who_id, whom_id) VALUES (?, ?)`
		db.Exec(query,
			userId, followsUserId,
		)

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodPost && unfollowsUsername != "" {
		fmt.Printf("\"/fllws/:username\" running as a Post-Method, where unfollow in c.FormParams()\n")

		unfollowsUserId, err := datalayer.GetUserId(unfollowsUsername, db)
		if err != nil {
			fmt.Printf("getUserId returned error: %v\n", err)
			return err
		}

		query := `DELETE FROM follower WHERE who_id=? and WHOM_ID=?`
		db.Exec(query,
			userId, unfollowsUserId,
		)

		return c.JSON(http.StatusNoContent, nil)

	} else if c.Request().Method == http.MethodGet {
		fmt.Printf("\"/fllws/:username\" running as a Get-Method\n")

		noFollowersStr := c.QueryParam("no")
		noFollowers := 100
		if noFollowersStr != "" {
			val, err := strconv.Atoi(noFollowersStr)
			if err == nil {
				noFollowers = val
			}
		}
		query := `SELECT user.username FROM user
                  INNER JOIN follower ON follower.whom_id=user.user_id
                  WHERE follower.who_id=?
                  LIMIT ?`

		rows, err := datalayer.QueryDB(db, query,
			userId, noFollowers,
		)
		if err != nil {
			fmt.Printf("messages: queryDB returned error: %v\n", err)
			return err
		}

		follows, err := helpers.RowsToMapList(rows)
		if err != nil {
			fmt.Printf("messages: rowsToMapList returned error: %v\n", err)
			return err
		}

		var followList []any

		for _, follow := range follows {
			followList = append(followList, follow["username"])
		}

		data := map[string]any{
			"follows": followList,
		}
		fmt.Printf("data: %v\n", data)

		return c.JSON(http.StatusOK, data)
	}

	fmt.Printf("ERROR: \"/fllws/:username\" was entered wrongly!\n")
	return c.JSON(http.StatusBadRequest, nil)
}

func FollowUser(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}

	row := db.QueryRow(`SELECT * FROM user
						WHERE username = ?`,
		username,
	)
	var user models.User
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.PwHash)
	if err != nil {
		fmt.Printf("row.Scan returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}
	db.Exec("insert into follower (who_id, whom_id) values (?, ?)", sessionUserId, user.UserID)
	err = helpers.AddFlash(c, fmt.Sprintf("You are now following \"%s\"", username))
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}

func UnfollowUser(c echo.Context) error {
	username := c.Param("username")
	fmt.Printf("User entered UserTimeline via route \"/:username\" as \"/%v\"\n", username)

	loggedIn, _ := helpers.IsUserLoggedIn(c)
	if !loggedIn {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}

	row := db.QueryRow(`SELECT * FROM user
						WHERE username = ?`,
		username,
	)
	var user models.User
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.PwHash)
	if err != nil {
		fmt.Printf("row.Scan returned error: %v\n", err)
		c.String(http.StatusNotFound, "Not found")
	}

	sessionUserId, err := helpers.GetSessionUserID(c)
	if err != nil {
		fmt.Printf("getSessionUserID returned error: %v\n", err)
		return err
	}
	db.Exec("delete from follower where who_id=? and whom_id=?", sessionUserId, user.UserID)

	err = helpers.AddFlash(c, fmt.Sprintf("You are no longer following \"%s\"", username))
	if err != nil {
		fmt.Printf("addFlash returned error: %v\n", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/%s", username))
}
