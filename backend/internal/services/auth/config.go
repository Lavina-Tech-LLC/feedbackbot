package auth

import (
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	ap := config.Confs.AuthProvider
	c.Data(lvn.Res(200, gin.H{
		"authorize_url": ap.BaseURL + "/authorize",
		"client_id":     ap.ClientID,
		"redirect_uri":  ap.RedirectURI,
		"scopes":        "openid profile email",
	}, ""))
}
