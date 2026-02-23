package svc_tenant_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/auth"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_tenant"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTenant_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "ct@example.com", "CT User")

	router := testutil.SetupRouter()
	router.POST("/tenants", auth.Auth, svc_tenant.CreateTenant)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, 0)
	body := map[string]string{"name": "New Org", "slug": "new-org"}
	w := testutil.DoRequest(router, "POST", "/tenants", body, token)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "New Org", data["name"])
	assert.Equal(t, "new-org", data["slug"])
}

func TestCreateTenant_NoAuth(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.POST("/tenants", auth.Auth, svc_tenant.CreateTenant)

	body := map[string]string{"name": "Org", "slug": "org"}
	w := testutil.DoRequest(router, "POST", "/tenants", body, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateTenant_MissingFields(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "mf@example.com", "MF User")

	router := testutil.SetupRouter()
	router.POST("/tenants", auth.Auth, svc_tenant.CreateTenant)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, 0)
	body := map[string]string{"name": "No Slug"}
	w := testutil.DoRequest(router, "POST", "/tenants", body, token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTenant_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "gt@example.com", "GT User")
	tenant := testutil.CreateTestTenant(t, user.ID, "Get Org", "get-org")

	router := testutil.SetupRouter()
	router.GET("/tenants/:id", auth.Auth, svc_tenant.GetTenant)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", "/tenants/"+uintToStr(tenant.ID), nil, token)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTenant_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "gtnf@example.com", "GTNF User")

	router := testutil.SetupRouter()
	router.GET("/tenants/:id", auth.Auth, svc_tenant.GetTenant)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, 0)
	w := testutil.DoRequest(router, "GET", "/tenants/99999", nil, token)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBots_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "bots@example.com", "Bots User")
	tenant := testutil.CreateTestTenant(t, user.ID, "Bot Org", "bot-org")

	router := testutil.SetupRouter()
	router.GET("/bots", auth.Auth, services.TenantMiddleware, svc_tenant.GetBots)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", "/bots", nil, token)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteBot_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "delbot@example.com", "Del User")
	tenant := testutil.CreateTestTenant(t, user.ID, "Del Org", "del-org")

	router := testutil.SetupRouter()
	router.DELETE("/bots/:id", auth.Auth, services.TenantMiddleware, svc_tenant.DeleteBot)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "DELETE", "/bots/99999", nil, token)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func uintToStr(n uint) string {
	return fmt.Sprintf("%d", n)
}
