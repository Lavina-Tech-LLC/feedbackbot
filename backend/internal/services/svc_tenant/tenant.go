package svc_tenant

import (
	"fmt"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type createTenantReq struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

func CreateTenant(c *gin.Context) {
	var req createTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "Invalid request: "+err.Error()))
		return
	}

	tenant := models.Tenant{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := models.DB.Create(&tenant).Error; err != nil {
		lvn.GinErr(c, 500, err, "Failed to create tenant")
		return
	}

	// Associate the user with the newly created tenant
	uid, exists := c.Get("user_id")
	if !exists {
		c.Data(lvn.Res(401, "", "user_id not found in context"))
		return
	}
	userID := fmt.Sprintf("%v", uid)

	ut := models.UserTenant{
		UserID:   userID,
		TenantID: tenant.ID,
	}
	if dbErr := models.DB.Where("user_id = ?", userID).FirstOrCreate(&ut).Error; dbErr != nil {
		lvn.GinErr(c, 500, dbErr, "Failed to associate user with tenant")
		return
	}

	c.Data(lvn.Res(201, tenant, ""))
}

func GetTenant(c *gin.Context) {
	id := c.Param("id")

	var tenant models.Tenant
	if err := models.DB.First(&tenant, id).Error; err != nil {
		c.Data(lvn.Res(404, "", "Tenant not found"))
		return
	}

	c.Data(lvn.Res(200, tenant, ""))
}
