package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware validates the JWT before allowing access to protected routes.
func JWTMiddleware(secret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		// Get the Authorization header sent by the client.
		// Expected format: "Bearer eyJhbGciOi..."
		authHeader := c.GetHeader("Authorization")

		// Reject the request if no Authorization header was provided.
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Split "Bearer <token>" into ["Bearer", "<token>"].
		parts := strings.SplitN(authHeader, " ", 2)

		// Reject anything that isn't in "Bearer <token>" format.
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header",
			})
			return
		}

		// Get the actual JWT string.
		tokenString := parts[1]

		// Create an empty Claims object.
		// ParseWithClaims will fill it with the data from the JWT.
		claims := &Claims{}

		// Parse the JWT, verify its signature, check its validity,
		// and put its payload into our Claims struct.
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,

			// This function tells the JWT library which secret
			// to use when verifying the token's signature.
			func(token *jwt.Token) (interface{}, error) {

				// Make sure the token uses an HMAC signing method.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				// Return our secret so the library can verify the signature.
				return []byte(secret), nil
			},
		)

		// Reject the request if the token is invalid, expired,
		// has a bad signature, etc.
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// The token is valid.
		// Store the authenticated user's information in Gin's context
		// so handlers/services can access it later.
		c.Set("user_id", claims.UserID)
		c.Set("client_id", claims.ClientID)

		// Allow the request to continue to the protected route.
		c.Next()
	}
}