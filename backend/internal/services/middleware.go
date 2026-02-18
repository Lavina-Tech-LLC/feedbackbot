package services

import (
	"strconv"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// TenantMiddleware extracts tenant_id from request header or query param
// and sets it in the Gin context. In the future, this will extract from auth token.
func TenantMiddleware(c *gin.Context) {
	// Try header first (X-Tenant-ID), then query param
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = c.Query("tenant_id")
	}

	if tenantIDStr == "" {
		c.Data(lvn.Res(400, "", "tenant_id is required"))
		c.Abort()
		return
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
	if err != nil {
		c.Data(lvn.Res(400, "", "invalid tenant_id"))
		c.Abort()
		return
	}

	c.Set("tenant_id", uint(tenantID))
	c.Next()
}

// GetTenantID extracts tenant_id from Gin context
func GetTenantID(c *gin.Context) uint {
	tenantID, _ := c.Get("tenant_id")
	if id, ok := tenantID.(uint); ok {
		return id
	}
	return 0
}
