package auth

import (
	"context"
	//"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coreos/go-oidc"
	// "golang.org/x/oauth2"
)

type KeycloakMiddleware struct {
	Verifier *oidc.IDTokenVerifier
}

func NewKeycloakMiddleware(issuer string, clientID string) (*KeycloakMiddleware, error) {
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	return &KeycloakMiddleware{
		Verifier: verifier,
	}, nil
}

// While integrating with keycloak Im facing with keycloak configration . aud is showing UUID of client-id instead of text value. below is the fix. make sure in keycloak. remove and added again the
// Client-scope
// https://chatgpt.com/share/67fe65df-3bdc-8007-96d5-01c5472d0570
func (km *KeycloakMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := km.Verifier.Verify(context.Background(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// You can store claims or roles in context if needed
		c.Set("user", claims)
		c.Next()
	}
}
