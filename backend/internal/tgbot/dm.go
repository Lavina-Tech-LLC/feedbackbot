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
		sendMessage(bot.Token, msg.Chat.ID, "👋 FeedbackBot ga xush kelibsiz!\n\nMenga xabar yuboring va men uni anonim ravishda guruh adminingizga yetkazaman.\n\nFaqat admin ko'rishi uchun xabaringiz oldidan /adminOnly buyrug'ini ishlating.")
		return
	}

	// Check if admin_only (case-insensitive prefix)
	adminOnly := false
	if len(text) >= len("/adminOnly") && strings.EqualFold(text[:len("/adminOnly")], "/adminOnly") {
		adminOnly = true
		text = strings.TrimSpace(text[len("/adminOnly"):])
		if text == "" {
			sendMessage(bot.Token, msg.Chat.ID, "Iltimos, /adminOnly buyrug'idan keyin fikringizni yozing.\n\nMisol: /adminOnly Men standup yig'ilishlarni yaxshilash kerak deb o'ylayman.")
			return
		}
	}

	if text == "" {
		sendMessage(bot.Token, msg.Chat.ID, "Iltimos, fikringizni matn ko'rinishida yuboring.")
		return
	}

	const maxFeedbackLen = 4000
	if len(text) > maxFeedbackLen {
		sendMessage(bot.Token, msg.Chat.ID, "Xabaringiz juda uzun. Iltimos, 4000 belgidan oshmasin.")
		return
	}

	// Find groups this user belongs to (via bot's tenant)
	var groups []models.Group
	models.DB.Where("bot_id = ? AND is_active = ?", bot.ID, true).Find(&groups)

	if len(groups) == 0 {
		sendMessage(bot.Token, msg.Chat.ID, "❌ Faol guruhlar topilmadi. Bot avval guruhga qo'shilishi kerak.")
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
	storePendingFeedback(bot.ID, userID, text, adminOnly)

	var keyboard [][]inlineButton
	for _, g := range groups {
		keyboard = append(keyboard, []inlineButton{
			{Text: g.Title, CallbackData: fmt.Sprintf("fb:%d", g.ID)},
		})
	}

	sendMessageWithKeyboard(bot.Token, msg.Chat.ID, "📋 Bu fikr-mulohaza qaysi guruh uchun?", keyboard)
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
				postText := fmt.Sprintf("📬 Anonim fikr-mulohaza:\n\n%s", message)
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
		sendMessage(bot.Token, chatID, "✅ Fikringiz adminga shaxsiy ravishda yuborildi. Guruhda e'lon qilinMAYDI.")
	} else {
		sendMessage(bot.Token, chatID, "✅ Fikringiz anonim ravishda yuborildi. Rahmat!")
	}
}

func storePendingFeedback(botID uint, userID int64, text string, adminOnly bool) {
	pf := models.PendingFeedback{
		TelegramUserID: userID,
		BotID:          botID,
		Text:           text,
		AdminOnly:      adminOnly,
	}
	var existing models.PendingFeedback
	if err := models.DB.Where("telegram_user_id = ?", userID).First(&existing).Error; err == nil {
		models.DB.Model(&existing).Updates(map[string]interface{}{
			"bot_id":     botID,
			"text":       text,
			"admin_only": adminOnly,
		})
	} else {
		models.DB.Create(&pf)
	}
}

func getPendingFeedback(userID int64) (models.PendingFeedback, bool) {
	var pf models.PendingFeedback
	if err := models.DB.Where("telegram_user_id = ?", userID).First(&pf).Error; err != nil {
		return pf, false
	}
	models.DB.Delete(&pf)
	return pf, true
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
