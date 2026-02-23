package tgbot

import (
	"testing"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(
		&models.Tenant{},
		&models.Bot{},
		&models.Group{},
		&models.FeedbackConfig{},
		&models.GroupUser{},
		&models.Feedback{},
		&models.User{},
		&models.UserTenant{},
		&models.PendingFeedback{},
	)
	models.DB = db
	config.Confs.Settings.JWTSecret = "test-secret"
}

func TestStorePendingFeedback(t *testing.T) {
	setupTestDB(t)

	storePendingFeedback(1, 12345, "hello", false)

	var pf models.PendingFeedback
	err := models.DB.Where("telegram_user_id = ?", 12345).First(&pf).Error
	assert.NoError(t, err)
	assert.Equal(t, "hello", pf.Text)
	assert.Equal(t, false, pf.AdminOnly)
	assert.Equal(t, uint(1), pf.BotID)
}

func TestStorePendingFeedback_Overwrites(t *testing.T) {
	setupTestDB(t)

	storePendingFeedback(1, 12345, "first message", false)
	storePendingFeedback(1, 12345, "second message", true)

	var pf models.PendingFeedback
	err := models.DB.Where("telegram_user_id = ?", 12345).First(&pf).Error
	assert.NoError(t, err)
	assert.Equal(t, "second message", pf.Text)
	assert.Equal(t, true, pf.AdminOnly)

	var count int64
	models.DB.Model(&models.PendingFeedback{}).Where("telegram_user_id = ?", 12345).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestGetPendingFeedback_Found(t *testing.T) {
	setupTestDB(t)

	storePendingFeedback(1, 12345, "pending msg", true)

	pf, ok := getPendingFeedback(12345)
	assert.True(t, ok)
	assert.Equal(t, "pending msg", pf.Text)
	assert.Equal(t, true, pf.AdminOnly)

	// Should be deleted after retrieval
	_, ok = getPendingFeedback(12345)
	assert.False(t, ok)
}

func TestGetPendingFeedback_NotFound(t *testing.T) {
	setupTestDB(t)

	_, ok := getPendingFeedback(99999)
	assert.False(t, ok)
}

func TestHandleMyChatMember_BotAdded(t *testing.T) {
	setupTestDB(t)

	tenant := models.Tenant{Name: "Test", Slug: "test"}
	models.DB.Create(&tenant)
	bot := models.Bot{TenantID: tenant.ID, Token: "tok", BotUsername: "bot", Verified: true}
	models.DB.Create(&bot)

	member := &ChatMemberUp{
		Chat: Chat{ID: -100123, Title: "Test Group", Type: "supergroup"},
		NewChatMember: ChatMember{
			Status: "member",
			User:   User{ID: 999, IsBot: true},
		},
	}

	handleMyChatMember(bot, member)

	var group models.Group
	err := models.DB.Where("chat_id = ?", -100123).First(&group).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Group", group.Title)
	assert.Equal(t, true, group.IsActive)
	assert.Equal(t, bot.ID, group.BotID)

	// FeedbackConfig should also be created
	var fc models.FeedbackConfig
	err = models.DB.Where("group_id = ?", group.ID).First(&fc).Error
	assert.NoError(t, err)
	assert.Equal(t, false, fc.PostToGroup)
}

func TestHandleMyChatMember_BotRemoved(t *testing.T) {
	setupTestDB(t)

	tenant := models.Tenant{Name: "Test2", Slug: "test2"}
	models.DB.Create(&tenant)
	bot := models.Bot{TenantID: tenant.ID, Token: "tok2", BotUsername: "bot2", Verified: true}
	models.DB.Create(&bot)
	group := models.Group{TenantID: tenant.ID, BotID: bot.ID, ChatID: -100456, Title: "Remove Group", Type: "supergroup", IsActive: true}
	models.DB.Create(&group)

	member := &ChatMemberUp{
		Chat: Chat{ID: -100456, Title: "Remove Group", Type: "supergroup"},
		NewChatMember: ChatMember{
			Status: "kicked",
			User:   User{ID: 999, IsBot: true},
		},
	}

	handleMyChatMember(bot, member)

	var updated models.Group
	models.DB.Where("chat_id = ?", -100456).First(&updated)
	assert.Equal(t, false, updated.IsActive)
}

func TestHandleMyChatMember_BotReAdded(t *testing.T) {
	setupTestDB(t)

	tenant := models.Tenant{Name: "Test3", Slug: "test3"}
	models.DB.Create(&tenant)
	bot := models.Bot{TenantID: tenant.ID, Token: "tok3", BotUsername: "bot3", Verified: true}
	models.DB.Create(&bot)
	group := models.Group{TenantID: tenant.ID, BotID: bot.ID, ChatID: -100789, Title: "Old Title", Type: "group", IsActive: false}
	models.DB.Create(&group)

	member := &ChatMemberUp{
		Chat: Chat{ID: -100789, Title: "New Title", Type: "supergroup"},
		NewChatMember: ChatMember{
			Status: "administrator",
			User:   User{ID: 999, IsBot: true},
		},
	}

	handleMyChatMember(bot, member)

	var updated models.Group
	models.DB.Where("chat_id = ?", -100789).First(&updated)
	assert.Equal(t, true, updated.IsActive)
	assert.Equal(t, "New Title", updated.Title)
	assert.Equal(t, "supergroup", updated.Type)
}

func TestHandleMyChatMember_PrivateChat_Ignored(t *testing.T) {
	setupTestDB(t)

	tenant := models.Tenant{Name: "Test4", Slug: "test4"}
	models.DB.Create(&tenant)
	bot := models.Bot{TenantID: tenant.ID, Token: "tok4", BotUsername: "bot4", Verified: true}
	models.DB.Create(&bot)

	member := &ChatMemberUp{
		Chat: Chat{ID: 12345, Title: "", Type: "private"},
		NewChatMember: ChatMember{
			Status: "member",
			User:   User{ID: 999, IsBot: true},
		},
	}

	handleMyChatMember(bot, member)

	var count int64
	models.DB.Model(&models.Group{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestSubmitFeedback_CreatesGroupUserAndFeedback(t *testing.T) {
	setupTestDB(t)

	tenant := models.Tenant{Name: "SF", Slug: "sf"}
	models.DB.Create(&tenant)
	bot := models.Bot{TenantID: tenant.ID, Token: "sf-tok", BotUsername: "sfbot", Verified: true}
	models.DB.Create(&bot)
	group := models.Group{TenantID: tenant.ID, BotID: bot.ID, ChatID: -200111, Title: "SF Group", Type: "supergroup", IsActive: true}
	models.DB.Create(&group)
	models.DB.Create(&models.FeedbackConfig{GroupID: group.ID, PostToGroup: false})

	// submitFeedback will try to send a Telegram message which will fail silently
	// We only verify DB state
	submitFeedback(bot, 12345, 67890, group, "test feedback", false)

	var gu models.GroupUser
	err := models.DB.Where("group_id = ? AND telegram_user_id = ?", group.ID, 67890).First(&gu).Error
	assert.NoError(t, err)

	var fb models.Feedback
	err = models.DB.Where("group_id = ? AND message = ?", group.ID, "test feedback").First(&fb).Error
	assert.NoError(t, err)
	assert.Equal(t, false, fb.AdminOnly)
	assert.Equal(t, gu.ID, fb.SenderID)
}

func TestSubmitFeedback_AdminOnly(t *testing.T) {
	setupTestDB(t)

	tenant := models.Tenant{Name: "AO", Slug: "ao"}
	models.DB.Create(&tenant)
	bot := models.Bot{TenantID: tenant.ID, Token: "ao-tok", BotUsername: "aobot", Verified: true}
	models.DB.Create(&bot)
	group := models.Group{TenantID: tenant.ID, BotID: bot.ID, ChatID: -200222, Title: "AO Group", Type: "supergroup", IsActive: true}
	models.DB.Create(&group)
	models.DB.Create(&models.FeedbackConfig{GroupID: group.ID, PostToGroup: true})

	submitFeedback(bot, 12345, 67890, group, "admin secret", true)

	var fb models.Feedback
	err := models.DB.Where("group_id = ? AND message = ?", group.ID, "admin secret").First(&fb).Error
	assert.NoError(t, err)
	assert.Equal(t, true, fb.AdminOnly)
	// Admin-only feedback should NOT be posted even if config says post_to_group
	assert.Equal(t, false, fb.Posted)
}
