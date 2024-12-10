package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"product_api/internal/core/domain"
	"time"
)

type TokenGenerator struct {
	secretKey []byte
}

func NewTokenGenerator(secretKey []byte) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
	}
}

func (g *TokenGenerator) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": domain.User{Username: "admin"},
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IDTokenCustomClaims holds structure of jwt claims of idToken
type IDTokenCustomClaims struct {
	User *domain.User `json:"user"`
	jwt.RegisteredClaims
}

func validateIDToken(tokenString string, key []byte) (*IDTokenCustomClaims, error) {
	claims := &IDTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*IDTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

func (g *TokenGenerator) ValidateIDToken(tokenString string) (*domain.User, error) {
	claims, err := validateIDToken(tokenString, g.secretKey) // uses public RSA key

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, fmt.Errorf("unauthorized")
	}

	return claims.User, nil
}
