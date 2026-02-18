package tgbot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
)

type Update struct {
	UpdateID     int64         `json:"update_id"`
	MyChatMember *ChatMemberUp `json:"my_chat_member"`
	Message      *Message      `json:"message"`
}

type ChatMemberUp struct {
	Chat          Chat       `json:"chat"`
	From          User       `json:"from"`
	NewChatMember ChatMember `json:"new_chat_member"`
	OldChatMember ChatMember `json:"old_chat_member"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type ChatMember struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Text      string `json:"text"`
}

type getUpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func StartPolling(bot models.Bot) {
	log.Printf("[tgbot] Starting polling for bot @%s (ID: %d)", bot.BotUsername, bot.ID)
	offset := int64(0)

	for {
		updates, err := getUpdates(bot.Token, offset)
		if err != nil {
			log.Printf("[tgbot] Error getting updates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			offset = update.UpdateID + 1
			handleUpdate(bot, update)
		}

		time.Sleep(1 * time.Second)
	}
}

func getUpdates(token string, offset int64) ([]Update, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=30&allowed_updates=[\"my_chat_member\",\"message\"]", token, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result getUpdatesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.Ok {
		return nil, fmt.Errorf("telegram API returned not ok")
	}

	return result.Result, nil
}

func handleUpdate(bot models.Bot, update Update) {
	if update.MyChatMember != nil {
		handleMyChatMember(bot, update.MyChatMember)
	}
}

func handleMyChatMember(bot models.Bot, member *ChatMemberUp) {
	chat := member.Chat
	newStatus := member.NewChatMember.Status

	// Only handle group/supergroup chats
	if chat.Type != "group" && chat.Type != "supergroup" {
		return
	}

	if newStatus == "member" || newStatus == "administrator" {
		// Bot added to group
		var group models.Group
		result := models.DB.Where("chat_id = ?", chat.ID).First(&group)
		if result.Error != nil {
			// Create new group
			group = models.Group{
				TenantID: bot.TenantID,
				BotID:    bot.ID,
				ChatID:   chat.ID,
				Title:    chat.Title,
				Type:     chat.Type,
				IsActive: true,
			}
			models.DB.Create(&group)

			// Create default feedback config
			config := models.FeedbackConfig{
				GroupID:     group.ID,
				PostToGroup: false,
			}
			models.DB.Create(&config)

			log.Printf("[tgbot] Bot added to group: %s (chat_id: %d)", chat.Title, chat.ID)
		} else {
			// Reactivate existing group
			models.DB.Model(&group).Updates(map[string]interface{}{
				"is_active": true,
				"title":     chat.Title,
				"type":      chat.Type,
			})
			log.Printf("[tgbot] Bot re-added to group: %s (chat_id: %d)", chat.Title, chat.ID)
		}
	} else if newStatus == "left" || newStatus == "kicked" {
		// Bot removed from group
		models.DB.Model(&models.Group{}).Where("chat_id = ?", chat.ID).Update("is_active", false)
		log.Printf("[tgbot] Bot removed from group: chat_id %d", chat.ID)
	}
}
