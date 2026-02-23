package svc_tenant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/tgbot"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type createBotReq struct {
	Token string `json:"token" binding:"required"`
}

type telegramGetMeResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID        int64  `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"result"`
	Description string `json:"description"`
}

func verifyBotToken(token string) (*telegramGetMeResponse, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result telegramGetMeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateBot(c *gin.Context) {
	var req createBotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "Invalid request: "+err.Error()))
		return
	}

	tenantID := services.GetTenantID(c)

	// Check for duplicate token
	var existingBot models.Bot
	if err := models.DB.Where("token = ?", req.Token).First(&existingBot).Error; err == nil {
		c.Data(lvn.Res(409, "", "Bot with this token already exists"))
		return
	}

	// Verify tenant exists
	var tenant models.Tenant
	if err := models.DB.First(&tenant, tenantID).Error; err != nil {
		c.Data(lvn.Res(404, "", "Tenant not found"))
		return
	}

	// Verify bot token with Telegram
	tgResp, err := verifyBotToken(req.Token)
	if err != nil {
		lvn.GinErr(c, 500, err, "Failed to verify bot token")
		return
	}

	if !tgResp.Ok {
		c.Data(lvn.Res(400, "", "Invalid bot token: "+tgResp.Description))
		return
	}

	bot := models.Bot{
		TenantID:    tenantID,
		Token:       req.Token,
		BotUsername: tgResp.Result.Username,
		BotName:     tgResp.Result.FirstName,
		Verified:    true,
	}

	if err := models.DB.Create(&bot).Error; err != nil {
		lvn.GinErr(c, 500, err, "Failed to create bot")
		return
	}

	go tgbot.StartPolling(bot)

	c.Data(lvn.Res(201, bot, ""))
}

func GetBots(c *gin.Context) {
	tenantID := services.GetTenantID(c)

	var bots []models.Bot
	models.DB.Scopes(db.TenantScope(tenantID)).Find(&bots)

	c.Data(lvn.Res(200, bots, ""))
}

func GetBot(c *gin.Context) {
	id := c.Param("id")
	tenantID := services.GetTenantID(c)

	var bot models.Bot
	if err := models.DB.Scopes(db.TenantScope(tenantID)).First(&bot, id).Error; err != nil {
		c.Data(lvn.Res(404, "", "Bot not found"))
		return
	}

	c.Data(lvn.Res(200, bot, ""))
}

func DeleteBot(c *gin.Context) {
	id := c.Param("id")
	tenantID := services.GetTenantID(c)

	var bot models.Bot
	if err := models.DB.Scopes(db.TenantScope(tenantID)).First(&bot, id).Error; err != nil {
		c.Data(lvn.Res(404, "", "Bot not found"))
		return
	}

	if err := models.DB.Delete(&bot).Error; err != nil {
		lvn.GinErr(c, 500, err, "Failed to delete bot")
		return
	}

	c.Data(lvn.Res(200, "", "Bot deleted"))
}
