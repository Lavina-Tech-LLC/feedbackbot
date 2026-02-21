package auth

import (
	"net/url"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	ap := config.Confs.AuthProvider
	authorizeURL := ap.BaseURL + "/api/auth/oauth/google/url?redirect_url=" + url.QueryEscape(ap.RedirectURI)
	c.Data(lvn.Res(200, gin.H{
		"authorize_url": authorizeURL,
		"client_id":     ap.ClientID,
		"redirect_uri":  ap.RedirectURI,
		"scopes":        "openid profile email",
	}, ""))
}
