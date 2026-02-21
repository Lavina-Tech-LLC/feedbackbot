package auth

import (
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// GetForgotPasswordURL returns the auth provider's password reset URL.
func GetForgotPasswordURL(c *gin.Context) {
	url := config.Confs.AuthProvider.BaseURL + "/forgot-password"
	c.Data(lvn.Res(200, gin.H{
		"url": url,
	}, ""))
}
