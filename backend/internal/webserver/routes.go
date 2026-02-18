package webServer

import (
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_tenant"
	"github.com/gin-gonic/gin"
)

func setRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	tenants := router.Group("/tenants")
	tenants.POST("", svc_tenant.CreateTenant)
	tenants.GET("/:id", svc_tenant.GetTenant)

	bots := router.Group("/bots")
	bots.POST("", svc_tenant.CreateBot)
	bots.GET("/:id", svc_tenant.GetBot)
	bots.DELETE("/:id", svc_tenant.DeleteBot)
}

func Listen() {
	router := gin.Default()
	setRoutes(router)

	addr := config.Confs.Settings.SrvAddress
	if addr == "" {
		addr = ":8080"
	}
	router.Run(addr)
}
