package user_repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"product_api/internal/core/domain"
)

type UserRepository struct {
	repo *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		repo: collection,
	}
}

func (u *UserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	cur, err := u.repo.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		var user domain.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.Background())

	return users, nil
}

func (u *UserRepository) Find(id string) (domain.User, error) {
	var user domain.User
	objectId, objectErr := primitive.ObjectIDFromHex(id)
	if objectErr != nil {
		return domain.User{}, objectErr
	}
	err := u.repo.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserRepository) Save(user domain.User) (domain.User, error) {

	_, err := u.repo.InsertOne(context.Background(), user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserRepository) Update(id string, user domain.User) (domain.User, error) {
	_, err := u.repo.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": user})
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserRepository) Delete(id string) error {
	objectId, objectErr := primitive.ObjectIDFromHex(id)
	if objectErr != nil {
		return objectErr
	}
	_, err := u.repo.DeleteOne(context.Background(), bson.M{"id": objectId})
	if err != nil {
		return err
	}
	return nil
}
