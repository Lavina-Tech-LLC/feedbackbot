package auth

import (
	"fmt"
	"strings"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth validates the Bearer JWT locally using the configured secret.
// Claims (user_id, email, role, name, tenant_id) are set in the Gin context for downstream handlers.
func Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.Data(lvn.Res(401, "", "Authorization required"))
		c.Abort()
		return
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Confs.Settings.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		c.Data(lvn.Res(401, "", "Invalid or expired token"))
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Data(lvn.Res(401, "", "Invalid token claims"))
		c.Abort()
		return
	}

	if v, ok := claims["user_id"]; ok {
		c.Set("user_id", v)
	}
	if v, ok := claims["email"]; ok {
		c.Set("email", v)
	}
	if v, ok := claims["name"]; ok {
		c.Set("name", v)
	}
	if v, ok := claims["role"]; ok {
		c.Set("role", v)
	}
	if v, ok := claims["tenant_id"]; ok {
		c.Set("tenant_id", v)
	}

	// Fallback: if tenant_id is not in the JWT, check the local DB
	if _, exists := c.Get("tenant_id"); !exists {
		if rawID, ok := c.Get("user_id"); ok {
			var userID uint
			switch v := rawID.(type) {
			case float64:
				userID = uint(v)
			}
			var ut models.UserTenant
			if err := models.DB.Where("user_id = ?", userID).First(&ut).Error; err == nil {
				c.Set("tenant_id", ut.TenantID)
			}
		}
	}

	c.Next()
}
