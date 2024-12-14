package services

import (
	"context"
	"crypto/rsa"
	"log"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
	"product_api/internal/utils"
)

type TokenService struct {
	tokenRepository       ports.ITokenRepository
	pub                   *rsa.PublicKey
	pri                   *rsa.PrivateKey
	refreshSecret         string
	iDExpirationSecs      int64
	refreshExpirationSecs int64
}

type TSConfig struct {
	TokenRepository       ports.ITokenRepository
	Pub                   *rsa.PublicKey
	Pri                   *rsa.PrivateKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) *TokenService {
	return &TokenService{
		tokenRepository:       c.TokenRepository,
		pub:                   c.Pub,
		pri:                   c.Pri,
		refreshSecret:         c.RefreshSecret,
		iDExpirationSecs:      c.IDExpirationSecs,
		refreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

func (t *TokenService) GenerateTokens(ctx context.Context, u *domain.User, prevAccessToken string) (*domain.Tokens, error) {
	accessToken, err := utils.CreateAccessToken(u, t.pri)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.CreateRefreshToken(u.ID, t.pri)
	if err != nil {
		return nil, err
	}

	if err := t.tokenRepository.SetRefreshToken(ctx, u.ID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.ID, err.Error())
		return nil, err
	}

	// delete user's current refresh token (used when refreshing idToken)
	if prevAccessToken != "" {
		if err := t.tokenRepository.DeleteRefreshToken(ctx, u.ID.String(), prevAccessToken); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.ID.String(), prevAccessToken)
		}
	}
	// store tokens

	return &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.SS,
	}, nil

}
