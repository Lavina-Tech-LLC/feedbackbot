package webServer

import (
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/auth"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_feedback"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_group"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services/svc_tenant"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authGroup := router.Group("/auth")
	authGroup.GET("/config", auth.GetConfig)
	authGroup.POST("/login", auth.Login)
	authGroup.POST("/register", auth.Register)
	authGroup.POST("/token", auth.ExchangeToken)
	authGroup.GET("/forgot-password", auth.GetForgotPasswordURL)
	authGroup.GET("/me", auth.Auth, auth.GetMe)

	tenants := router.Group("/tenants")
	tenants.POST("", svc_tenant.CreateTenant)
	tenants.GET("/:id", svc_tenant.GetTenant)

	bots := router.Group("/bots", auth.Auth, services.TenantMiddleware)
	bots.POST("", svc_tenant.CreateBot)
	bots.GET("/:id", svc_tenant.GetBot)
	bots.DELETE("/:id", svc_tenant.DeleteBot)

	groups := router.Group("/groups", auth.Auth, services.TenantMiddleware)
	groups.GET("", svc_group.GetGroups)
	groups.GET("/:id", svc_group.GetGroup)
	groups.PATCH("/:id", svc_group.UpdateGroup)
	groups.PATCH("/:id/config", svc_group.UpdateGroupConfig)

	feedbacks := router.Group("/feedbacks", auth.Auth, services.TenantMiddleware)
	feedbacks.GET("", svc_feedback.GetFeedbacks)
}

func Listen() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type", "X-Tenant-ID"},
	}))

	setRoutes(router)

	addr := config.Confs.Settings.SrvAddress
	if addr == "" {
		addr = ":8080"
	}
	router.Run(addr)
}
