package auth

import (
	"strings"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth validates the JWT Bearer token from the Authorization header.
// Claims (user_id, email, role) are set in the Gin context for downstream handlers.
func Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.Data(lvn.Res(401, "", "Authorization required"))
		c.Abort()
		return
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Confs.JWT.AccessSecret), nil
	})

	if err != nil || !token.Valid {
		c.Data(lvn.Res(401, "", "Invalid or expired token"))
		c.Abort()
		return
	}

	if userID, ok := claims["user_id"]; ok {
		c.Set("user_id", userID)
	}
	if email, ok := claims["email"]; ok {
		c.Set("email", email)
	}
	if role, ok := claims["role"]; ok {
		c.Set("role", role)
	}

	c.Next()
}
