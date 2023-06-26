package helpers

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// authSecrete ...
var authSecrete = GetEnv("AUTH_SECRET_KEY", "secrete")

// Token ...
type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// AuthClaims is the struct that will be encoded to a JWT.
type AuthClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateAuthTokens generates both the detailed token and refresh token
// tokenExpiry is optional and defaults to 3 hours for access token and 7 days for refresh token
func GenerateAuthTokens(AuthClaims *AuthClaims, tokenExpiry ...time.Duration) (*Token, error) {
	accessTokenExpiry := time.Now().Add(time.Hour * 3).Unix()       // 3 hours
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7).Unix() // 7 days
	if len(tokenExpiry) > 0 {
		accessTokenExpiry = time.Now().Add(tokenExpiry[0]).Unix()
	}

	if len(tokenExpiry) > 1 {
		refreshTokenExpiry = time.Now().Add(tokenExpiry[1]).Unix()
	}

	AuthClaims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: accessTokenExpiry,
		Issuer:    "ochom",
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims).SignedString([]byte(authSecrete))
	if err != nil {
		return nil, err
	}

	AuthClaims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: refreshTokenExpiry,
		Issuer:    "ochom",
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims).SignedString([]byte(authSecrete))
	if err != nil {
		return nil, err
	}

	return &Token{AccessToken: token, RefreshToken: refreshToken}, nil
}

// ValidateToken validates the token
func ValidateToken(token string) (*AuthClaims, error) {
	AuthClaims := &AuthClaims{}
	tkn, err := jwt.ParseWithClaims(token, AuthClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(authSecrete), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}
	return AuthClaims, nil
}
