package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// Login proxies a login request to the auth provider and returns JWT tokens.
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "email and password are required"))
		return
	}

	resp, err := forwardToAuthProvider("/login", req)
	if err != nil {
		c.Data(lvn.Res(502, "", "auth provider unavailable"))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.Data(lvn.Res(resp.StatusCode, "", extractError(body)))
		return
	}

	data, err := parseAuthResponse(body)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to parse auth response"))
		return
	}

	ensureTenant(data)

	c.Data(lvn.Res(200, data, ""))
}

// Register proxies a registration request to the auth provider and returns JWT tokens.
func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "name, email, and password are required"))
		return
	}

	resp, err := forwardToAuthProvider("/register", req)
	if err != nil {
		c.Data(lvn.Res(502, "", "auth provider unavailable"))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.Data(lvn.Res(resp.StatusCode, "", extractError(body)))
		return
	}

	data, err := parseAuthResponse(body)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to parse auth response"))
		return
	}

	ensureTenant(data)

	c.Data(lvn.Res(200, data, ""))
}

// forwardToAuthProvider sends a JSON POST request to the auth provider.
func forwardToAuthProvider(path string, payload interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := config.Confs.AuthProvider.BaseURL + path
	return http.Post(url, "application/json", bytes.NewReader(jsonBody))
}

// parseAuthResponse parses the auth provider JSON response and extracts claims from the access token.
func parseAuthResponse(body []byte) (gin.H, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	// The auth provider may nest tokens under "data" or return them at the top level.
	data, ok := raw["data"].(map[string]interface{})
	if !ok {
		data = raw
	}

	accessToken, _ := data["access_token"].(string)
	refreshToken, _ := data["refresh_token"].(string)

	result := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	// Extract claims from the access token (without verification — the auth provider already signed it).
	if accessToken != "" {
		claims := jwt.MapClaims{}
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
		if t, _, err := parser.ParseUnverified(accessToken, claims); err == nil && t != nil {
			if v, ok := claims["user_id"]; ok {
				result["user_id"] = v
			}
			if v, ok := claims["email"]; ok {
				result["email"] = v
			}
			if v, ok := claims["tenant_id"]; ok {
				result["tenant_id"] = v
			}
			if v, ok := claims["name"]; ok {
				result["name"] = v
			}
		}
	}

	return result, nil
}

// ensureTenant auto-creates a tenant in the DB if the tenant_id from the token doesn't exist yet.
func ensureTenant(data gin.H) {
	tidRaw, ok := data["tenant_id"]
	if !ok {
		return
	}

	var tenantID uint
	switch v := tidRaw.(type) {
	case float64:
		tenantID = uint(v)
	case uint:
		tenantID = v
	}
	if tenantID == 0 {
		return
	}

	var tenant models.Tenant
	if err := models.DB.First(&tenant, tenantID).Error; err != nil {
		emailStr := fmt.Sprintf("%v", data["email"])
		tenant = newTenantFromEmail(emailStr)
		models.DB.Create(&tenant)
	}
}

// extractError tries to pull an error message from the auth provider JSON response.
func extractError(body []byte) string {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return "authentication failed"
	}
	if msg, ok := raw["message"].(string); ok && msg != "" {
		return msg
	}
	if msg, ok := raw["error"].(string); ok && msg != "" {
		return msg
	}
	return "authentication failed"
}
