package services

import (
	"context"
	"crypto/rsa"
	"product_api/internal/core/domain"
	"product_api/internal/utils"
)

type TokenService struct {
	pub           *rsa.PublicKey
	pri           *rsa.PrivateKey
	RefreshSecret string
}

func NewTokenService(pub *rsa.PublicKey, pri *rsa.PrivateKey, refreshSecret string) *TokenService {
	return &TokenService{
		pub:           pub,
		pri:           pri,
		RefreshSecret: refreshSecret,
	}
}

func (t *TokenService) GenerateTokens(ctx context.Context, u *domain.User, prevAccessToken string) (*domain.Tokens, error) {
	accessToken, err := utils.CreateAccessToken(u, t.pri)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.CreateRefreshToken(u.ID, t.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.SS,
	}, nil

}
