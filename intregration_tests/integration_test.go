package intregrationtests

import (
	"context"
	"flag"
	"log"
	"os"
	"product_api/internal/adapters/repository/sareerepo"
	"product_api/internal/core/domain"
	"product_api/internal/core/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// All tests that use mtest.Setup() are expected to be integration tests, so skip them when the
	// -short flag is included in the "go test" command. Also, we have to parse flags here to use
	// testing.Short() because flags aren't parsed before TestMain() is called.
	flag.Parse()
	if testing.Short() {
		log.Print("skipping mtest integration test in short mode")
		return
	}

	if err := mtest.Setup(); err != nil {
		log.Fatal(err)
	}
	defer os.Exit(m.Run())
	if err := mtest.Teardown(); err != nil {
		log.Fatal(err)
	}
}

func Test_sareeRepository(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mtest.ClusterURI()))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("test")
	repo := db.Collection("sarees")

	t.Run("Should get all sarees from find all", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		repo.InsertMany(context.Background(), []interface{}{
			bson.D{
				{Key: "fabrictype", Value: "cotton"},
				{Key: "color", Value: "red"},
				{Key: "category", Value: "saree"},
			},
			bson.D{
				{Key: "color", Value: "blue"},
				{Key: "category", Value: "saree"},
				{Key: "fabrictype", Value: "cotton"},
			},
			bson.D{
				{Key: "fabrictype", Value: "cotton"},
				{Key: "color", Value: "red"},
				{Key: "category", Value: "saree"},
			},
		})

		sareeRepo := sareerepo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		sarees, err := sareeService.FindAll()
		require.NoError(t, err)
		assert.Equal(t, 3, len(sarees))
		assert.Equal(t, "cotton", sarees[0].FabricType)
		assert.Equal(t, "red", sarees[0].Color)
		assert.Equal(t, "saree", sarees[0].Category)
	})

	t.Run("Should get saree by id", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		res, err := repo.InsertOne(context.Background(), bson.D{
			{Key: "fabrictype", Value: "cotton"},
			{Key: "color", Value: "red"},
			{Key: "category", Value: "saree"},
		})
		require.NoError(t, err)

		sareeRepo := sareerepo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		saree, err := sareeService.Find(res.InsertedID.(primitive.ObjectID).Hex())
		require.NoError(t, err)
		assert.Equal(t, "cotton", saree.FabricType)
		assert.Equal(t, "red", saree.Color)
		assert.Equal(t, "saree", saree.Category)
	})

	t.Run("Should save saree", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		sareeRepo := sareerepo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		saree, err := sareeService.Save(domain.Saree{
			FabricType: "cotton",
			Color:      "red",
			Category:   "saree",
		})
		require.NoError(t, err)
		assert.Equal(t, "cotton", saree.FabricType)
		assert.Equal(t, "red", saree.Color)
		assert.Equal(t, "saree", saree.Category)
	})

	t.Run("Should update saree", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		res, err := repo.InsertOne(context.Background(), bson.D{
			{Key: "fabrictype", Value: "cotton"},
			{Key: "color", Value: "red"},
			{Key: "category", Value: "saree"},
		})
		require.NoError(t, err)

		sareeRepo := sareerepo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		saree, err := sareeService.Update(res.InsertedID.(primitive.ObjectID).Hex(), domain.Saree{
			FabricType: "cotton",
			Color:      "red",
			Category:   "saree",
		})
		require.NoError(t, err)
		assert.Equal(t, "cotton", saree.FabricType)
		assert.Equal(t, "red", saree.Color)
		assert.Equal(t, "saree", saree.Category)
	})

	t.Run("Should delete saree", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		res, err := repo.InsertOne(context.Background(), bson.D{
			{Key: "fabrictype", Value: "cotton"},
			{Key: "color", Value: "red"},
			{Key: "category", Value: "saree"},
		})
		require.NoError(t, err)

		sareeRepo := sareerepo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		err = sareeService.Delete(res.InsertedID.(primitive.ObjectID).Hex())
		require.NoError(t, err)
	})

}
