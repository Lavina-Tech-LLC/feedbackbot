package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const TestJWTSecret = "test-secret-key"

// SetupTestDB initializes an in-memory SQLite database with all models migrated.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(
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
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	models.DB = db
	config.Confs.Settings.JWTSecret = TestJWTSecret
	return db
}

// CreateTestUser creates a user and returns it.
func CreateTestUser(t *testing.T, email, name string) models.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.MinCost)
	user := models.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := models.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

// CreateTestTenant creates a tenant and associates it with a user.
func CreateTestTenant(t *testing.T, userID uint, name, slug string) models.Tenant {
	t.Helper()
	tenant := models.Tenant{Name: name, Slug: slug}
	if err := models.DB.Create(&tenant).Error; err != nil {
		t.Fatalf("failed to create test tenant: %v", err)
	}
	ut := models.UserTenant{UserID: userID, TenantID: tenant.ID}
	if err := models.DB.Create(&ut).Error; err != nil {
		t.Fatalf("failed to create user-tenant link: %v", err)
	}
	return tenant
}

// GenerateTestToken generates a JWT access token for testing.
func GenerateTestToken(userID uint, email, name, role string, tenantID uint) string {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"email":     email,
		"name":      name,
		"role":      role,
		"tenant_id": tenantID,
		"type":      "access",
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(TestJWTSecret))
	return token
}

// GenerateTestRefreshToken generates a JWT refresh token for testing.
func GenerateTestRefreshToken(userID uint) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(TestJWTSecret))
	return token
}

// GenerateExpiredToken generates an expired JWT token.
func GenerateExpiredToken(userID uint) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(TestJWTSecret))
	return token
}

// SetupRouter creates a Gin engine in test mode.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// DoRequest performs an HTTP request against the given handler/router.
func DoRequest(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
