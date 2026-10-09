package saree_repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"product_api/internal/core/domain"

	"github.com/google/uuid"
)

type SareeRepository struct {
	repo *mongo.Collection
}

func NewSareeRepository(collection *mongo.Collection) *SareeRepository {
	return &SareeRepository{
		repo: collection,
	}
}

// Sarees are keyed by their UID field, which the driver stores as "uid".
func byUID(id string) (bson.M, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return bson.M{"uid": uid}, nil
}

// BackfillUIDs gives every saree saved without an ID (older records had an
// empty one) its own, so it can be edited and deleted. Returns how many it fixed.
func (s *SareeRepository) BackfillUIDs(ctx context.Context) (int, error) {
	filter := bson.M{"$or": bson.A{
		bson.M{"uid": bson.M{"$exists": false}},
		bson.M{"uid": uuid.Nil},
	}}
	cur, err := s.repo.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return 0, err
	}
	var docs []struct {
		ID interface{} `bson:"_id"`
	}
	if err := cur.All(ctx, &docs); err != nil {
		return 0, err
	}
	for _, doc := range docs {
		if _, err := s.repo.UpdateByID(ctx, doc.ID, bson.M{"$set": bson.M{"uid": uuid.New()}}); err != nil {
			return 0, err
		}
	}
	return len(docs), nil
}

func (s *SareeRepository) FindAll() ([]domain.Saree, error) {
	// Non-nil so an empty catalog encodes as [] rather than null.
	sarees := []domain.Saree{}
	// Oldest first, so the storefront order matches the order sarees were added.
	cur, err := s.repo.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var saree domain.Saree
		err := cur.Decode(&saree)
		if err != nil {
			return nil, err
		}
		sarees = append(sarees, saree)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return sarees, nil
}

func (s *SareeRepository) Find(id string) (domain.Saree, error) {
	var saree domain.Saree
	filter, err := byUID(id)
	if err != nil {
		return domain.Saree{}, err
	}
	err = s.repo.FindOne(context.Background(), filter).Decode(&saree)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Saree{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeRepository) Save(saree domain.Saree) (domain.Saree, error) {

	_, err := s.repo.InsertOne(context.Background(), saree)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeRepository) Update(id string, saree domain.Saree) (domain.Saree, error) {
	filter, err := byUID(id)
	if err != nil {
		return domain.Saree{}, err
	}
	result, err := s.repo.UpdateOne(context.Background(), filter, bson.M{"$set": saree})
	if err != nil {
		return domain.Saree{}, err
	}
	if result.MatchedCount == 0 {
		return domain.Saree{}, domain.ErrNotFound
	}
	return saree, nil

}

func (s *SareeRepository) Delete(id string) error {
	filter, err := byUID(id)
	if err != nil {
		return err
	}
	result, err := s.repo.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
