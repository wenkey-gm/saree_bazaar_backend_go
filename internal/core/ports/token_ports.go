package ports

import (
	"context"
	"product_api/internal/core/domain"
	"time"
)

type ITokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}

type ITokenService interface {
	GenerateTokens(ctx context.Context, u *domain.User, prevAccessToken string) (*domain.Tokens, error)
	ValidateIDToken(tokenString string) (*domain.User, error)
	SignOut(ctx context.Context, userID string) error
}
