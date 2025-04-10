package repo_wrappers

import (
	"context"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/metrics"
	"minitwit/src/models"
	"minitwit/src/utils"

	"github.com/labstack/echo/v4"
)

func CreateMessage(c echo.Context, authorID int, text string) error {
	newMessage := helpers.NewMessage(authorID, text)
	err := messageRepo.Create(c.Request().Context(), newMessage)
	if err != nil {
		log.Printf("messageRepo.Create returned error: %v\n", err)
		return err
	}

	pubTime := utils.GetTimeFromInt(newMessage.PubDate)
	metrics.MessagesPosts.WithLabelValues(
		utils.GetHourAsString(pubTime), 
		utils.GetWeekdayAsString(pubTime),
		).Inc()
	return nil
}

func GetMessagesFiltered(c echo.Context, conditions map[string]any, noMsgs int) ([]models.Message, error) {
	msgs, err := messageRepo.GetFiltered(c.Request().Context(), conditions, noMsgs, "pub_date DESC")
	if err != nil {
		log.Printf("GetMessagesFiltered: messageRepo.GetFiltered returned error: %v\n", err)
		return nil, err
	}
	return msgs, nil
}

func EnhanceMessages(c echo.Context, msgs []models.Message, isAPI bool) ([]map[string]any) {
	var enrichedMsgs []map[string]any
	for _, msg := range msgs {
		enrichedMsg := map[string]any{
			"pub_date": msg.PubDate,
		}

		author, _ := GetUserByID(c, msg.AuthorID)
		
		var username, email string
		if author != nil {
			username = author.Username
			email = author.Email
		} else {
			log.Printf("⚠️ Warning: Could not find user for message author_id=%d\n", msg.AuthorID)
			username = "Unknown"
			email = ""
		}

		if isAPI {
			enrichedMsg["content"] = msg.Text
			enrichedMsg["user"] = username
		} else {
			enrichedMsg["text"] = msg.Text
			enrichedMsg["username"] = username
			enrichedMsg["email"] = email
		}

		enrichedMsgs = append(enrichedMsgs, enrichedMsg)
	}
	return enrichedMsgs
}

func CountAllMessages(c context.Context) (int, error) {
	count, err := messageRepo.CountAll(c)
	if err != nil {
		log.Printf("❌ Repository Error: CountAllMessages returned error: %v\n", err)
		return 0, err
	}
	return count, nil
}

func CountFilteredMessages(c context.Context, conditions map[string]any) (int, error) {
	msgs, err := messageRepo.GetFiltered(c, conditions, 0, "")
	if err != nil {
		log.Printf("❌ Repository Error: CountFilteredMessages returned error: %v\n", err)
		return 0, err
	}
	return len(msgs), nil
}