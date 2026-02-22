package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// Auth validates the Bearer token by forwarding it to the auth provider.
// Claims (user_id, email, role, name, tenant_id) are set in the Gin context for downstream handlers.
func Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.Data(lvn.Res(401, "", "Authorization required"))
		c.Abort()
		return
	}

	req, err := http.NewRequest("GET", config.Confs.AuthProvider.BaseURL+"/api/user/me", nil)
	if err != nil {
		c.Data(lvn.Res(401, "", "Invalid or expired token"))
		c.Abort()
		return
	}
	req.Header.Set("Authorization", header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		c.Data(lvn.Res(401, "", "Invalid or expired token"))
		c.Abort()
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.Data(lvn.Res(401, "", "Invalid or expired token"))
		c.Abort()
		return
	}

	// Auth provider returns {"user": {...}} format
	data, _ := result["user"].(map[string]interface{})
	if data == nil {
		data, _ = result["data"].(map[string]interface{})
	}
	if data == nil {
		data = result
	}

	// Map "id" → "user_id" for consistency
	if v, ok := data["id"]; ok {
		c.Set("user_id", v)
	}
	for _, key := range []string{"email", "role", "name"} {
		if val, ok := data[key]; ok {
			c.Set(key, val)
		}
	}

	// Extract tenant_id from the JWT claims (auth provider doesn't return it in /api/user/me)
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	if parts := strings.Split(tokenStr, "."); len(parts) == 3 {
		if payload, err := base64DecodeSegment(parts[1]); err == nil {
			var claims map[string]interface{}
			if json.Unmarshal(payload, &claims) == nil {
				if v, ok := claims["tenant_id"]; ok {
					c.Set("tenant_id", v)
				}
			}
		}
	}

	// Fallback: if tenant_id is not in the JWT, check the local DB
	if _, exists := c.Get("tenant_id"); !exists {
		userID := fmt.Sprintf("%v", c.MustGet("user_id"))
		var ut models.UserTenant
		if err := models.DB.Where("user_id = ?", userID).First(&ut).Error; err == nil {
			c.Set("tenant_id", ut.TenantID)
		}
	}

	c.Next()
}

// base64DecodeSegment decodes a JWT base64url segment.
func base64DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(seg)
}
