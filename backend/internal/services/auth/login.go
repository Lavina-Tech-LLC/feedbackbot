package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register creates a new user with bcrypt-hashed password and returns JWT tokens.
func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "name, email, and password are required"))
		return
	}

	// Check if user already exists
	var existing models.User
	if err := models.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.Data(lvn.Res(409, "", "email already registered"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to hash password"))
		return
	}

	user := models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := models.DB.Create(&user).Error; err != nil {
		c.Data(lvn.Res(500, "", "failed to create user"))
		return
	}

	// Auto-create tenant from email domain
	tenant := newTenantFromEmail(user.Email)
	models.DB.Create(&tenant)
	models.DB.Create(&models.UserTenant{
		UserID:   fmt.Sprintf("%d", user.ID),
		TenantID: tenant.ID,
	})

	tokens, err := generateTokens(user, tenant.ID)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to generate tokens"))
		return
	}

	c.Data(lvn.Res(200, tokens, ""))
}

// Login finds a user by email, verifies the bcrypt password, and returns JWT tokens.
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "email and password are required"))
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.Data(lvn.Res(401, "", "invalid email or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.Data(lvn.Res(401, "", "invalid email or password"))
		return
	}

	// Look up tenant
	var tenantID uint
	var ut models.UserTenant
	if err := models.DB.Where("user_id = ?", fmt.Sprintf("%d", user.ID)).First(&ut).Error; err == nil {
		tenantID = ut.TenantID
	}

	tokens, err := generateTokens(user, tenantID)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to generate tokens"))
		return
	}

	c.Data(lvn.Res(200, tokens, ""))
}

// RefreshToken validates a refresh token and issues new access + refresh tokens.
func RefreshToken(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "refresh_token is required"))
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Confs.Settings.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		c.Data(lvn.Res(401, "", "invalid or expired refresh token"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Data(lvn.Res(401, "", "invalid token claims"))
		return
	}

	if claims["type"] != "refresh" {
		c.Data(lvn.Res(401, "", "invalid token type"))
		return
	}

	// Look up user by ID from claims
	userIDRaw, ok := claims["user_id"]
	if !ok {
		c.Data(lvn.Res(401, "", "invalid token claims"))
		return
	}

	var userID uint
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.Data(lvn.Res(401, "", "user not found"))
		return
	}

	var tenantID uint
	var ut models.UserTenant
	if err := models.DB.Where("user_id = ?", fmt.Sprintf("%d", user.ID)).First(&ut).Error; err == nil {
		tenantID = ut.TenantID
	}

	tokens, err := generateTokens(user, tenantID)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to generate tokens"))
		return
	}

	c.Data(lvn.Res(200, tokens, ""))
}

func generateTokens(user models.User, tenantID uint) (gin.H, error) {
	secret := []byte(config.Confs.Settings.JWTSecret)

	accessClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"name":      user.Name,
		"role":      user.Role,
		"tenant_id": tenantID,
		"type":      "access",
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_id":       user.ID,
		"email":         user.Email,
		"name":          user.Name,
		"role":          user.Role,
		"tenant_id":     tenantID,
	}, nil
}

func newTenantFromEmail(email string) models.Tenant {
	domain := "unknown"
	if parts := strings.SplitN(email, "@", 2); len(parts) == 2 {
		domain = parts[1]
	}
	return models.Tenant{
		Name: domain,
		Slug: domain,
	}
}
