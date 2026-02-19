package services

import (
	"strconv"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// TenantMiddleware extracts tenant_id from JWT claims (set by Auth middleware),
// falling back to X-Tenant-ID header or query param for backward compatibility.
func TenantMiddleware(c *gin.Context) {
	// Try JWT context first (set by Auth middleware)
	if tid, exists := c.Get("tenant_id"); exists {
		switch v := tid.(type) {
		case float64:
			c.Set("tenant_id", uint(v))
			c.Next()
			return
		case uint:
			c.Next()
			return
		}
	}

	// Fall back to header / query param
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
