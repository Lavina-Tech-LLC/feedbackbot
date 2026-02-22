package auth

import (
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")
	name, _ := c.Get("name")
	tenantID, _ := c.Get("tenant_id")

	c.Data(lvn.Res(200, gin.H{
		"user_id":   userID,
		"email":     email,
		"name":      name,
		"tenant_id": tenantID,
		"role":      role,
	}, ""))
}
