package auth

import (
	"fmt"
	"strings"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")
	name, _ := c.Get("name")

	var tenantID uint
	if tid, ok := c.Get("tenant_id"); ok {
		switch v := tid.(type) {
		case float64:
			tenantID = uint(v)
		case uint:
			tenantID = v
		}
	}

	// Look up or auto-create tenant
	if tenantID > 0 {
		var tenant models.Tenant
		if err := models.DB.First(&tenant, tenantID).Error; err != nil {
			// Tenant not found — create one from email domain
			tenant = newTenantFromEmail(fmt.Sprintf("%v", email))
			if err := models.DB.Create(&tenant).Error; err != nil {
				c.Data(lvn.Res(500, "", "failed to create tenant"))
				return
			}
			tenantID = tenant.ID
		}
	}

	c.Data(lvn.Res(200, gin.H{
		"user_id":   userID,
		"email":     email,
		"name":      name,
		"tenant_id": tenantID,
		"role":      role,
	}, ""))
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
