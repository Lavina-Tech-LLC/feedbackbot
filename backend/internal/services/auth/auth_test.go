package auth_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/auth"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.POST("/auth/register", auth.Register)

	body := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	w := testutil.DoRequest(router, "POST", "/auth/register", body, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	assert.Equal(t, "test@example.com", data["email"])
}

func TestRegister_DuplicateEmail(t *testing.T) {
	testutil.SetupTestDB(t)
	testutil.CreateTestUser(t, "dup@example.com", "Dup User")

	router := testutil.SetupRouter()
	router.POST("/auth/register", auth.Register)

	body := map[string]string{
		"name":     "Another User",
		"email":    "dup@example.com",
		"password": "password123",
	}
	w := testutil.DoRequest(router, "POST", "/auth/register", body, "")

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_MissingFields(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.POST("/auth/register", auth.Register)

	body := map[string]string{"email": "test@example.com"}
	w := testutil.DoRequest(router, "POST", "/auth/register", body, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "login@example.com", "Login User")
	testutil.CreateTestTenant(t, user.ID, "Test Org", "test-org")

	router := testutil.SetupRouter()
	router.POST("/auth/login", auth.Login)

	body := map[string]string{
		"email":    "login@example.com",
		"password": "testpassword",
	}
	w := testutil.DoRequest(router, "POST", "/auth/login", body, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	testutil.SetupTestDB(t)
	testutil.CreateTestUser(t, "user@example.com", "User")

	router := testutil.SetupRouter()
	router.POST("/auth/login", auth.Login)

	body := map[string]string{
		"email":    "user@example.com",
		"password": "wrongpassword",
	}
	w := testutil.DoRequest(router, "POST", "/auth/login", body, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_NonExistentUser(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.POST("/auth/login", auth.Login)

	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	w := testutil.DoRequest(router, "POST", "/auth/login", body, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_MissingFields(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.POST("/auth/login", auth.Login)

	body := map[string]string{"email": "test@example.com"}
	w := testutil.DoRequest(router, "POST", "/auth/login", body, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefreshToken_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "refresh@example.com", "Refresh User")
	testutil.CreateTestTenant(t, user.ID, "Test Org", "test-org-r")

	router := testutil.SetupRouter()
	router.POST("/auth/refresh", auth.RefreshToken)

	refreshToken := testutil.GenerateTestRefreshToken(user.ID)
	body := map[string]string{"refresh_token": refreshToken}
	w := testutil.DoRequest(router, "POST", "/auth/refresh", body, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
}

func TestRefreshToken_AccessTokenRejected(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "at@example.com", "AT User")

	router := testutil.SetupRouter()
	router.POST("/auth/refresh", auth.RefreshToken)

	accessToken := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, 0)
	body := map[string]string{"refresh_token": accessToken}
	w := testutil.DoRequest(router, "POST", "/auth/refresh", body, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.POST("/auth/refresh", auth.RefreshToken)

	body := map[string]string{"refresh_token": "not-a-valid-jwt"}
	w := testutil.DoRequest(router, "POST", "/auth/refresh", body, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMe_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "me@example.com", "Me User")
	tenant := testutil.CreateTestTenant(t, user.ID, "Me Org", "me-org")

	router := testutil.SetupRouter()
	router.GET("/auth/me", auth.Auth, auth.GetMe)

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", "/auth/me", nil, token)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "me@example.com", data["email"])
}

func TestGetMe_NoAuth(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()
	router.GET("/auth/me", auth.Auth, auth.GetMe)

	w := testutil.DoRequest(router, "GET", "/auth/me", nil, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMe_ExpiredToken(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "exp@example.com", "Exp User")

	router := testutil.SetupRouter()
	router.GET("/auth/me", auth.Auth, auth.GetMe)

	token := testutil.GenerateExpiredToken(user.ID)
	w := testutil.DoRequest(router, "GET", "/auth/me", nil, token)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTenantMiddleware_FromJWT(t *testing.T) {
	testutil.SetupTestDB(t)
	user := testutil.CreateTestUser(t, "tenant@example.com", "Tenant User")
	tenant := testutil.CreateTestTenant(t, user.ID, "TM Org", "tm-org")

	router := testutil.SetupRouter()
	router.GET("/test", auth.Auth, services.TenantMiddleware, func(c *gin.Context) {
		tid := services.GetTenantID(c)
		c.JSON(200, gin.H{"tenant_id": tid})
	})

	token := testutil.GenerateTestToken(user.ID, user.Email, user.Name, user.Role, tenant.ID)
	w := testutil.DoRequest(router, "GET", "/test", nil, token)

	assert.Equal(t, http.StatusOK, w.Code)
}
