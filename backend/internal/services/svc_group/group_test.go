package svc_group_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/auth"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_group"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestGroup(t *testing.T, tenantID, botID uint, chatID int64, title string) models.Group {
	t.Helper()
	group := models.Group{
		TenantID: tenantID,
		BotID:    botID,
		ChatID:   chatID,
		Title:    title,
		Type:     "supergroup",
		IsActive: true,
	}
	if err := models.DB.Create(&group).Error; err != nil {
		t.Fatalf("failed to create test group: %v", err)
	}
	config := models.FeedbackConfig{
		GroupID:     group.ID,
		PostToGroup: false,
	}
	models.DB.Create(&config)
	return group
}

func createTestBot(t *testing.T, tenantID uint) models.Bot {
	t.Helper()
	bot := models.Bot{
		TenantID:    tenantID,
		Token:       "test-token-" + fmt.Sprintf("%d", tenantID),
		BotUsername: "testbot",
		BotName:     "Test Bot",
		Verified:    true,
	}
	if err := models.DB.Create(&bot).Error; err != nil {
		t.Fatalf("failed to create test bot: %v", err)
	}
	return bot
}

func TestGetGroups_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "grp@example.com", "Grp User")
	tenant := testutil.CreateTestTenant(t, user.ID, "Grp Org", "grp-org")
	bot := createTestBot(t, tenant.ID)
	createTestGroup(t, tenant.ID, bot.ID, -1001234567890, "Test Group")

	router := testutil.SetupRouter()
	router.GET("/groups", auth.Auth, services.TenantMiddleware, svc_group.GetGroups)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", "/groups", nil, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].([]interface{})
	assert.Len(t, data, 1)
}

func TestGetGroups_IsolatedByTenant(t *testing.T) {
	testutil.SetupTestDB(t)
	user1 := testutil.CreateTestUser(t, "u1@example.com", "User1")
	tenant1 := testutil.CreateTestTenant(t, user1.ID, "Org1", "org1")
	bot1 := createTestBot(t, tenant1.ID)
	createTestGroup(t, tenant1.ID, bot1.ID, -100111, "Group A")

	user2 := testutil.CreateTestUser(t, "u2@example.com", "User2")
	tenant2 := testutil.CreateTestTenant(t, user2.ID, "Org2", "org2")
	bot2 := createTestBot(t, tenant2.ID)
	createTestGroup(t, tenant2.ID, bot2.ID, -100222, "Group B")

	router := testutil.SetupRouter()
	router.GET("/groups", auth.Auth, services.TenantMiddleware, svc_group.GetGroups)

	// User1 should only see Group A
	token1 := testutil.GenerateTestToken(user1.ID, user1.Email, user1.Name, user1.Role, tenant1.ID)
	w1 := testutil.DoRequest(router, "GET", "/groups", nil, token1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var resp1 map[string]interface{}
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &resp1))
	data1 := resp1["data"].([]interface{})
	assert.Len(t, data1, 1)
}

func TestGetGroup_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "gg@example.com", "GG User")
	tenant := testutil.CreateTestTenant(t, user.ID, "GG Org", "gg-org")
	bot := createTestBot(t, tenant.ID)
	group := createTestGroup(t, tenant.ID, bot.ID, -100333, "GG Group")

	router := testutil.SetupRouter()
	router.GET("/groups/:id", auth.Auth, services.TenantMiddleware, svc_group.GetGroup)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", fmt.Sprintf("/groups/%d", group.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetGroup_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "ggnf@example.com", "GGNF User")
	tenant := testutil.CreateTestTenant(t, user.ID, "GGNF Org", "ggnf-org")

	router := testutil.SetupRouter()
	router.GET("/groups/:id", auth.Auth, services.TenantMiddleware, svc_group.GetGroup)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", "/groups/99999", nil, token)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateGroup_ToggleActive(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "ug@example.com", "UG User")
	tenant := testutil.CreateTestTenant(t, user.ID, "UG Org", "ug-org")
	bot := createTestBot(t, tenant.ID)
	group := createTestGroup(t, tenant.ID, bot.ID, -100444, "UG Group")

	router := testutil.SetupRouter()
	router.PATCH("/groups/:id", auth.Auth, services.TenantMiddleware, svc_group.UpdateGroup)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	body := map[string]interface{}{"is_active": false}
	w := testutil.DoRequest(router, "PATCH", fmt.Sprintf("/groups/%d", group.ID), body, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, false, data["is_active"])
}

func TestUpdateGroupConfig_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "ugc@example.com", "UGC User")
	tenant := testutil.CreateTestTenant(t, user.ID, "UGC Org", "ugc-org")
	bot := createTestBot(t, tenant.ID)
	group := createTestGroup(t, tenant.ID, bot.ID, -100555, "UGC Group")

	router := testutil.SetupRouter()
	router.PATCH("/groups/:id/config", auth.Auth, services.TenantMiddleware, svc_group.UpdateGroupConfig)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	body := map[string]interface{}{"post_to_group": true}
	w := testutil.DoRequest(router, "PATCH", fmt.Sprintf("/groups/%d/config", group.ID), body, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, true, data["post_to_group"])
}
