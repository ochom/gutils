// Package auth provides JWT token generation and validation utilities.
//
// This package offers a simple way to generate and validate JWT tokens for authentication
// purposes. It supports generating both access tokens and refresh tokens with configurable
// expiration times.
//
// Example usage:
//
//	// Generate authentication tokens for a user
//	userData := map[string]string{
//		"user_id": "12345",
//		"email":   "user@example.com",
//		"role":    "admin",
//	}
//
//	tokens, err := auth.GenerateAuthTokens(userData)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Access Token:", tokens["token"])
//	fmt.Println("Refresh Token:", tokens["refreshToken"])
//
//	// Validate and extract claims from a token
//	claims, err := auth.GetAuthClaims(tokens["token"])
//	if err != nil {
//		log.Fatal("Invalid token")
//	}
//	fmt.Println("User ID:", claims["user_id"])
package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var authSecrete string = "secrete"

// authClaims is the struct that will be encoded to a JWT.
type authClaims struct {
	Data map[string]string `json:"data"`
	jwt.StandardClaims
}

// GenerateAuthTokens generates both an access token and a refresh token.
//
// The data map is embedded in both tokens and can contain user information such as
// user ID, email, role, etc.
//
// By default, the access token expires in 3 hours and the refresh token in 7 days.
// Custom expiration times can be provided as optional parameters:
//   - tokenExpiry[0]: access token expiration duration
//   - tokenExpiry[1]: refresh token expiration duration
//
// Returns a map containing:
//   - "token": the access token
//   - "refreshToken": the refresh token
//
// Example:
//
//	// Default expiration times
//	tokens, err := auth.GenerateAuthTokens(map[string]string{"user_id": "123"})
//
//	// Custom expiration: 1 hour access, 24 hours refresh
//	tokens, err := auth.GenerateAuthTokens(
//		map[string]string{"user_id": "123"},
//		time.Hour,
//		24*time.Hour,
//	)
func GenerateAuthTokens(data map[string]string, tokenExpiry ...time.Duration) (map[string]string, error) {
	accessTokenExpiry := time.Now().Add(time.Hour * 3).Unix()       // 3 hours
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7).Unix() // 7 days
	if len(tokenExpiry) > 0 {
		accessTokenExpiry = time.Now().Add(tokenExpiry[0]).Unix()
	}

	if len(tokenExpiry) > 1 {
		refreshTokenExpiry = time.Now().Add(tokenExpiry[1]).Unix()
	}

	claims := &authClaims{Data: data}
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: accessTokenExpiry,
		Issuer:    "ochom",
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(authSecrete))
	if err != nil {
		return nil, err
	}

	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: refreshTokenExpiry,
		Issuer:    "ochom",
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(authSecrete))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token":        token,
		"refreshToken": refreshToken,
	}, nil
}

// GetAuthClaims validates a JWT token and extracts the embedded data claims.
//
// Returns the data map that was originally passed to GenerateAuthTokens,
// or an error if the token is invalid, expired, or malformed.
//
// Example:
//
//	claims, err := auth.GetAuthClaims(tokenString)
//	if err != nil {
//		// Token is invalid or expired
//		return unauthorized()
//	}
//
//	userID := claims["user_id"]
//	role := claims["role"]
func GetAuthClaims(token string) (map[string]string, error) {
	claims := &authClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(authSecrete), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	return claims.Data, nil
}
