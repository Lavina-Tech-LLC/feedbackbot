package auth

import (
	"strings"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth is a Gin middleware that validates JWT Bearer tokens.
// It extracts user_id from claims and sets it in the context.
func Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.Data(lvn.Res(401, "", "unauthorized"))
		c.Abort()
		return
	}

	tokenStr, found := strings.CutPrefix(header, "Bearer ")
	if !found {
		c.Data(lvn.Res(401, "", "unauthorized"))
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Confs.JWT.AccessSecret), nil
	})
	if err != nil || !token.Valid {
		c.Data(lvn.Res(401, "", "unauthorized"))
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Data(lvn.Res(401, "", "unauthorized"))
		c.Abort()
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.Data(lvn.Res(401, "", "unauthorized"))
		c.Abort()
		return
	}

	c.Set("user_id", uint(userIDFloat))
	c.Next()
}
