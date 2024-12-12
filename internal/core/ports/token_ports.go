package ports

import (
	"context"
	"product_api/internal/core/domain"
)

type ITokenService interface {
	GenerateTokens(ctx context.Context, u *domain.User, prevAccessToken string) (*domain.Tokens, error)
}
