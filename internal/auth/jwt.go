package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the information we want to store inside our JWT.
type Claims struct {
	UserID   int `json:"user_id"`
	ClientID int `json:"client_id"`

	// Adds standard JWT fields such as expiration (exp).
	jwt.RegisteredClaims
}

// GenerateToken creates and signs a JWT after successful login.
func GenerateToken(userID int, clientID int, secret string) (string, error) {

	// Create the data that will be stored inside the token.
	claims := Claims{
		UserID:   userID,
		ClientID: clientID,

		// Set standard JWT fields.
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expires 24 hours from now.
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	// Create the JWT using HMAC-SHA256.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign the token with our secret and return the resulting string.
	return token.SignedString([]byte(secret))
}