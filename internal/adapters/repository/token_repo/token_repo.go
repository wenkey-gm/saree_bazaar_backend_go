package token_repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TokenRepository struct {
	repo *mongo.Collection
}

func NewTokenRepository(collection *mongo.Collection) *TokenRepository {
	return &TokenRepository{
		repo: collection,
	}
}

func (t *TokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	_, err := t.repo.InsertOne(ctx, bson.M{"user_id": userID, "token_id": tokenID, "expires_in": expiresIn})
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenRepository) DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error {
	_, err := t.repo.DeleteOne(ctx, bson.M{"user_id": userID, "token_id": prevTokenID})
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	_, err := t.repo.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	return nil
}
