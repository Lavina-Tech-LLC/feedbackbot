package svc_tenant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
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
	userID := fmt.Sprintf("%v", c.MustGet("user_id"))
	authHeader := c.GetHeader("Authorization")

	// Try to update tenant_id on the auth provider; never block tenant creation.
	if err := patchAuthProviderTenantID(c.Request.Context(), authHeader, tenant.ID); err != nil {
		// Auth provider doesn't support PATCH or failed — store locally
		ut := models.UserTenant{
			UserID:   userID,
			TenantID: tenant.ID,
		}
		if dbErr := models.DB.Create(&ut).Error; dbErr != nil {
			lvn.GinErr(c, 500, dbErr, "Failed to associate user with tenant")
			return
		}
	}

	c.Data(lvn.Res(201, tenant, ""))
}

// patchAuthProviderTenantID tries to PATCH the user's tenant_id on the auth provider.
// Uses a 5-second timeout so a slow/unresponsive auth provider cannot hang tenant creation.
func patchAuthProviderTenantID(parent context.Context, authHeader string, tenantID uint) error {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	body, _ := json.Marshal(map[string]interface{}{"tenant_id": tenantID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, config.Confs.AuthProvider.BaseURL+"/api/user/me", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("auth provider returned status %d", resp.StatusCode)
	}
	return nil
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
