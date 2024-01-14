package helpers

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// getAuthSecrete ...
func getAuthSecrete() string {
	return GetEnv("AUTH_SECRET_KEY", "secrete")
}

// authClaims is the struct that will be encoded to a JWT.
type authClaims struct {
	Data map[string]string `json:"data"`
	jwt.StandardClaims
}

// GenerateAuthTokens generates both the detailed token and refresh token
// tokenExpiry is optional and defaults to 3 hours for access token and 7 days for refresh token
func GenerateAuthTokens(data map[string]string, tokenExpiry ...time.Duration) (string, string, error) {
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

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(getAuthSecrete()))
	if err != nil {
		return "", "", err
	}

	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: refreshTokenExpiry,
		Issuer:    "ochom",
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(getAuthSecrete()))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// GetAuthClaims ...
func GetAuthClaims(token string) (map[string]string, error) {
	claims := &authClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(getAuthSecrete()), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	return claims.Data, nil
}
