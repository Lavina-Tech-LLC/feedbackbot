package tgbot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
)

func handlePrivateMessage(bot models.Bot, msg *Message) {
	userID := msg.From.ID
	text := strings.TrimSpace(msg.Text)

	if text == "/start" {
		sendMessage(bot.Token, msg.Chat.ID, "👋 Welcome to FeedbackBot!\n\nSend me a message and I'll deliver it anonymously to your team admin.\n\nUse /adminOnly before your message to keep it visible only to the admin.")
		return
	}

	// Check if admin_only (case-insensitive prefix)
	adminOnly := false
	if len(text) >= len("/adminOnly") && strings.EqualFold(text[:len("/adminOnly")], "/adminOnly") {
		adminOnly = true
		text = strings.TrimSpace(text[len("/adminOnly"):])
		if text == "" {
			sendMessage(bot.Token, msg.Chat.ID, "Please write your feedback after /adminOnly.\n\nExample: /adminOnly I think we should improve our standup meetings.")
			return
		}
	}

	if text == "" {
		sendMessage(bot.Token, msg.Chat.ID, "Please send a text message with your feedback.")
		return
	}

	// Find groups this user belongs to (via bot's tenant)
	var groups []models.Group
	models.DB.Where("bot_id = ? AND is_active = ?", bot.ID, true).Find(&groups)

	if len(groups) == 0 {
		sendMessage(bot.Token, msg.Chat.ID, "❌ No active groups found. The bot needs to be added to a group first.")
		return
	}

	if len(groups) == 1 {
		// Auto-assign to the only group
		submitFeedback(bot, msg.Chat.ID, userID, groups[0], text, adminOnly)
		return
	}

	// Multiple groups — store pending feedback and show keyboard
	// For now, use the first group (TODO: implement inline keyboard picker in future iteration)
	// Store in a simple way: use callback data pattern
	storePendingFeedback(userID, text, adminOnly)

	var keyboard [][]inlineButton
	for _, g := range groups {
		keyboard = append(keyboard, []inlineButton{
			{Text: g.Title, CallbackData: fmt.Sprintf("fb:%d", g.ID)},
		})
	}

	sendMessageWithKeyboard(bot.Token, msg.Chat.ID, "📋 Which group is this feedback for?", keyboard)
}

func submitFeedback(bot models.Bot, chatID int64, telegramUserID int64, group models.Group, message string, adminOnly bool) {
	// Find or create GroupUser
	var groupUser models.GroupUser
	result := models.DB.Where("group_id = ? AND telegram_user_id = ?", group.ID, telegramUserID).First(&groupUser)
	if result.Error != nil {
		groupUser = models.GroupUser{
			TenantID:       group.TenantID,
			GroupID:        group.ID,
			TelegramUserID: telegramUserID,
		}
		models.DB.Create(&groupUser)
	}

	// Create feedback
	feedback := models.Feedback{
		TenantID:  group.TenantID,
		GroupID:   group.ID,
		SenderID:  groupUser.ID,
		Message:   message,
		AdminOnly: adminOnly,
		Posted:    false,
	}
	models.DB.Create(&feedback)

	// Post to group if config allows and not admin_only
	if !adminOnly {
		var config models.FeedbackConfig
		if err := models.DB.Where("group_id = ?", group.ID).First(&config).Error; err == nil {
			if config.PostToGroup {
				postText := fmt.Sprintf("📬 Anonymous Feedback:\n\n%s", message)
				if config.ForumTopicID != nil && *config.ForumTopicID > 0 {
					sendMessageToTopic(bot.Token, group.ChatID, *config.ForumTopicID, postText)
				} else {
					sendMessage(bot.Token, group.ChatID, postText)
				}
				models.DB.Model(&feedback).Update("posted", true)
			}
		}
	}

	// Confirm to user
	if adminOnly {
		sendMessage(bot.Token, chatID, "✅ Your feedback has been sent privately to the admin. It will NOT be posted in the group.")
	} else {
		sendMessage(bot.Token, chatID, "✅ Your feedback has been submitted anonymously. Thank you!")
	}
}

// Simple in-memory pending feedback store
var pendingFeedback = make(map[int64]pendingFB)

type pendingFB struct {
	Text      string
	AdminOnly bool
}

func storePendingFeedback(userID int64, text string, adminOnly bool) {
	pendingFeedback[userID] = pendingFB{Text: text, AdminOnly: adminOnly}
}

func getPendingFeedback(userID int64) (pendingFB, bool) {
	fb, ok := pendingFeedback[userID]
	if ok {
		delete(pendingFeedback, userID)
	}
	return fb, ok
}

type inlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func sendMessage(token string, chatID int64, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {fmt.Sprintf("%d", chatID)},
		"text":    {text},
	})
	if err != nil {
		log.Printf("[tgbot] Error sending message: %v", err)
		return
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
}

func sendMessageToTopic(token string, chatID int64, topicID int, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":           {fmt.Sprintf("%d", chatID)},
		"text":              {text},
		"message_thread_id": {fmt.Sprintf("%d", topicID)},
	})
	if err != nil {
		log.Printf("[tgbot] Error sending message to topic: %v", err)
		return
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
}

func sendMessageWithKeyboard(token string, chatID int64, text string, keyboard [][]inlineButton) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	kbJSON, _ := json.Marshal(map[string]interface{}{
		"inline_keyboard": keyboard,
	})

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":      {fmt.Sprintf("%d", chatID)},
		"text":         {text},
		"reply_markup": {string(kbJSON)},
	})
	if err != nil {
		log.Printf("[tgbot] Error sending keyboard message: %v", err)
		return
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
}
