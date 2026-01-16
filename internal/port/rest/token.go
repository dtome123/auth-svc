package rest

import (
	"auth-svc/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	GrantType           string `json:"grant_type"`
	ClientAssertionType string `json:"client_assertion_type"`
	ClientAssertion     string `json:"client_assertion"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (s *RestServer) Token(c *gin.Context) {

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	res, err := s.svc.GetAuthService().Token(c.Request.Context(), auth.TokenInput{
		GrantType:           req.GrantType,
		ClientAssertionType: req.ClientAssertionType,
		ClientAssertion:     req.ClientAssertion,
	})

	if err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, TokenResponse{
		AccessToken: res.AccessToken,
		TokenType:   res.TokenType,
		ExpiresIn:   res.ExpiresIn,
	})
}
