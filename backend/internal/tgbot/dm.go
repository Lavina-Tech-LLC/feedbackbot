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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[tgbot] PANIC in handlePrivateMessage: %v", r)
			sendMessage(bot.Token, msg.Chat.ID, "❌ An unexpected error occurred. Please try again.")
		}
	}()

	userID := msg.From.ID
	text := strings.TrimSpace(msg.Text)
	lang := detectLang(msg.From.LanguageCode)

	// Commands
	if text == "/start" || text == "/help" {
		sendMessage(bot.Token, msg.Chat.ID, getMsg("welcome", lang))
		return
	}

	// Unknown commands (starts with / but not known)
	if strings.HasPrefix(text, "/") && !strings.HasPrefix(text, "/adminOnly") {
		sendMessage(bot.Token, msg.Chat.ID, getMsg("unknownCommand", lang))
		return
	}

	// Check if admin_only
	adminOnly := false
	if strings.HasPrefix(text, "/adminOnly") {
		adminOnly = true
		text = strings.TrimSpace(strings.TrimPrefix(text, "/adminOnly"))
		if text == "" {
			sendMessage(bot.Token, msg.Chat.ID, getMsg("adminOnlyEmpty", lang))
			return
		}
	}

	if text == "" {
		sendMessage(bot.Token, msg.Chat.ID, getMsg("emptyMessage", lang))
		return
	}

	// Rate limiting
	if !checkRateLimit(userID) {
		sendMessage(bot.Token, msg.Chat.ID, getMsg("rateLimited", lang))
		return
	}

	// Find groups this user belongs to (via bot's tenant)
	var groups []models.Group
	models.DB.Where("bot_id = ? AND is_active = ?", bot.ID, true).Find(&groups)

	if len(groups) == 0 {
		sendMessage(bot.Token, msg.Chat.ID, getMsg("noGroups", lang))
		return
	}

	if len(groups) == 1 {
		submitFeedback(bot, msg.Chat.ID, userID, groups[0], text, adminOnly, lang)
		return
	}

	// Multiple groups — store pending feedback and show keyboard
	storePendingFeedback(userID, text, adminOnly, lang)

	var keyboard [][]inlineButton
	for _, g := range groups {
		keyboard = append(keyboard, []inlineButton{
			{Text: g.Title, CallbackData: fmt.Sprintf("fb:%d", g.ID)},
		})
	}

	sendMessageWithKeyboard(bot.Token, msg.Chat.ID, getMsg("pickGroup", lang), keyboard)
}

func submitFeedback(bot models.Bot, chatID int64, telegramUserID int64, group models.Group, message string, adminOnly bool, lang string) {
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
		sendMessage(bot.Token, chatID, getMsg("feedbackSentAdminOnly", lang))
	} else {
		sendMessage(bot.Token, chatID, getMsg("feedbackSent", lang))
	}
}

// Simple in-memory pending feedback store
var pendingFeedback = make(map[int64]pendingFB)

type pendingFB struct {
	Text      string
	AdminOnly bool
	Lang      string
}

func storePendingFeedback(userID int64, text string, adminOnly bool, lang string) {
	pendingFeedback[userID] = pendingFB{Text: text, AdminOnly: adminOnly, Lang: lang}
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
		"chat_id":    {fmt.Sprintf("%d", chatID)},
		"text":       {text},
		"parse_mode": {"Markdown"},
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
