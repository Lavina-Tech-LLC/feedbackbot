package auth

import (
	"io"
	"net/http"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type tokenExchangeRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
}

// ExchangeToken takes an OAuth code, forwards it to the auth provider's OAuth callback
// endpoint, and returns JWT tokens along with user info.
func ExchangeToken(c *gin.Context) {
	var req tokenExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "code is required"))
		return
	}

	resp, err := forwardToAuthProvider("/api/auth/oauth/google/callback", gin.H{"code": req.Code})
	if err != nil {
		c.Data(lvn.Res(502, "", "auth provider unavailable"))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.Data(lvn.Res(resp.StatusCode, "", extractError(body)))
		return
	}

	data, err := parseAuthResponse(body)
	if err != nil {
		c.Data(lvn.Res(500, "", "failed to parse auth response"))
		return
	}

	ensureTenant(data)

	c.Data(lvn.Res(200, data, ""))
}
