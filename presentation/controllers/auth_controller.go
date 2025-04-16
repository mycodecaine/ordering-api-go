package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// TokenRequest represents the request body to get a Keycloak token
type TokenRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret"` // optional
	Realm        string `json:"realm" binding:"required"`
}

// TokenResponse maps Keycloak token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// AuthController handles token generation
type AuthController struct {
	KeycloakBaseURL string
}

func NewAuthController(keycloakBaseURL string) *AuthController {
	return &AuthController{KeycloakBaseURL: keycloakBaseURL}
}

// GetToken godoc
// @Summary Generate Keycloak token
// @Description Authenticate a user and get a Keycloak access token
// @Tags auth
// @Accept json
// @Produce json
// @Param tokenRequest body TokenRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/token [post]
func (ac *AuthController) GetToken(c *gin.Context) {
	var req TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build token URL and form data
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", ac.KeycloakBaseURL, req.Realm)
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", req.ClientID)
	data.Set("username", req.Username)
	data.Set("password", req.Password)
	if req.ClientSecret != "" {
		data.Set("client_secret", req.ClientSecret)
	}

	// Make request to Keycloak
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact Keycloak"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		c.JSON(resp.StatusCode, gin.H{"error": errorBody})
		return
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	c.JSON(http.StatusOK, token)
}
