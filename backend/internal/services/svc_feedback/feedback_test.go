package svc_feedback_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/auth"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_feedback"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupFeedbackTestData(t *testing.T) (models.User, models.Tenant, models.Group) {
	t.Helper()
	user := testutil.CreateTestUser(t, "fb@example.com", "FB User")
	tenant := testutil.CreateTestTenant(t, user.ID, "FB Org", "fb-org")

	bot := models.Bot{
		TenantID:    tenant.ID,
		Token:       "test-bot-token",
		BotUsername: "fbbot",
		BotName:     "FB Bot",
		Verified:    true,
	}
	models.DB.Create(&bot)

	group := models.Group{
		TenantID: tenant.ID,
		BotID:    bot.ID,
		ChatID:   -100999,
		Title:    "FB Group",
		Type:     "supergroup",
		IsActive: true,
	}
	models.DB.Create(&group)

	groupUser := models.GroupUser{
		TenantID:       tenant.ID,
		GroupID:        group.ID,
		TelegramUserID: 12345,
	}
	models.DB.Create(&groupUser)

	// Create some feedbacks
	feedbacks := []models.Feedback{
		{TenantID: tenant.ID, GroupID: group.ID, SenderID: groupUser.ID, Message: "Public feedback 1", AdminOnly: false},
		{TenantID: tenant.ID, GroupID: group.ID, SenderID: groupUser.ID, Message: "Admin only feedback", AdminOnly: true},
		{TenantID: tenant.ID, GroupID: group.ID, SenderID: groupUser.ID, Message: "Public feedback 2", AdminOnly: false},
	}
	for _, fb := range feedbacks {
		models.DB.Create(&fb)
	}

	return user, tenant, group
}

func TestGetFeedbacks_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user, tenant, group := setupFeedbackTestData(t)

	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/feedbacks?group_id=%d", group.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	feedbacks := data["data"].([]interface{})
	assert.Len(t, feedbacks, 3)
	assert.Equal(t, float64(3), data["total"])
}

func TestGetFeedbacks_FilterAdminOnly(t *testing.T) {
	testutil.SetupTestDB(t)
	user, tenant, group := setupFeedbackTestData(t)

	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/feedbacks?group_id=%d&admin_only=true", group.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	feedbacks := data["data"].([]interface{})
	assert.Len(t, feedbacks, 1)
}

func TestGetFeedbacks_FilterPublicOnly(t *testing.T) {
	testutil.SetupTestDB(t)
	user, tenant, group := setupFeedbackTestData(t)

	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/feedbacks?group_id=%d&admin_only=false", group.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	feedbacks := data["data"].([]interface{})
	assert.Len(t, feedbacks, 2)
}

func TestGetFeedbacks_Pagination(t *testing.T) {
	testutil.SetupTestDB(t)
	user, tenant, group := setupFeedbackTestData(t)

	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/feedbacks?group_id=%d&page=1&limit=2", group.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	feedbacks := data["data"].([]interface{})
	assert.Len(t, feedbacks, 2)
	assert.Equal(t, float64(3), data["total"])
}

func TestGetFeedbacks_SenderIDNotExposed(t *testing.T) {
	testutil.SetupTestDB(t)
	user, tenant, group := setupFeedbackTestData(t)

	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/feedbacks?group_id=%d&limit=1", group.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)
	// SenderID has json:"-" so it should not appear
	assert.NotContains(t, w.Body.String(), "sender_id")
}

func TestGetFeedbacks_TenantIsolation(t *testing.T) {
	testutil.SetupTestDB(t)

	// Create tenant 1 with feedbacks
	user1 := testutil.CreateTestUser(t, "iso1@example.com", "Iso1")
	tenant1 := testutil.CreateTestTenant(t, user1.ID, "Iso1 Org", "iso1-org")
	bot1 := models.Bot{TenantID: tenant1.ID, Token: "tok1", BotUsername: "bot1", Verified: true}
	models.DB.Create(&bot1)
	group1 := models.Group{TenantID: tenant1.ID, BotID: bot1.ID, ChatID: -100001, Title: "G1", Type: "supergroup", IsActive: true}
	models.DB.Create(&group1)
	gu1 := models.GroupUser{TenantID: tenant1.ID, GroupID: group1.ID, TelegramUserID: 111}
	models.DB.Create(&gu1)
	models.DB.Create(&models.Feedback{TenantID: tenant1.ID, GroupID: group1.ID, SenderID: gu1.ID, Message: "T1 feedback"})

	// Create tenant 2
	user2 := testutil.CreateTestUser(t, "iso2@example.com", "Iso2")
	tenant2 := testutil.CreateTestTenant(t, user2.ID, "Iso2 Org", "iso2-org")

	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	// Tenant 2 should see no feedbacks
	token2 := testutil.GenerateTestToken(user2.ID, user2.Email, user2.Name, user2.Role, tenant2.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/feedbacks?group_id=%d", group1.ID), nil, token2)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	feedbacks := data["data"].([]interface{})
	assert.Len(t, feedbacks, 0)
}

func TestGetFeedbacks_NoAuth(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.GET("/feedbacks", auth.Auth, services.TenantMiddleware, svc_feedback.GetFeedbacks)

	w := testutil.DoRequest(router, "GET", "/feedbacks", nil, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
