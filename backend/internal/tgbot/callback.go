package tgbot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
)

func handleCallbackQuery(bot models.Bot, cq *CallbackQuery) {
	// Answer callback to remove loading state
	answerCallback(bot.Token, cq.ID)

	if !strings.HasPrefix(cq.Data, "fb:") {
		return
	}

	groupIDStr := strings.TrimPrefix(cq.Data, "fb:")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		log.Printf("[tgbot] Invalid callback data: %s", cq.Data)
		return
	}

	// Get pending feedback
	pending, ok := getPendingFeedback(cq.From.ID)
	if !ok {
		sendMessage(bot.Token, cq.Message.Chat.ID, "⏳ Session expired. Please send your feedback again.")
		return
	}

	// Find group
	var group models.Group
	if err := models.DB.First(&group, groupID).Error; err != nil {
		sendMessage(bot.Token, cq.Message.Chat.ID, "❌ Group not found.")
		return
	}

	submitFeedback(bot, cq.Message.Chat.ID, cq.From.ID, group, pending.Text, pending.AdminOnly)
}

func answerCallback(token string, callbackID string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token)
	resp, err := http.PostForm(apiURL, url.Values{
		"callback_query_id": {callbackID},
	})
	if err != nil {
		return
	}
	resp.Body.Close()
}
