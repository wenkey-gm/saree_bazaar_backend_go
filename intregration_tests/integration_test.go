package intregrationtests

import (
	"context"
	"flag"
	"log"
	"os"
	"product_api/internal/adapters/repository/saree_repo"
	"product_api/internal/core/domain"
	"product_api/internal/core/services"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
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

		sareeRepo := saree_repo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		sarees, err := sareeService.FindAll()
		require.NoError(t, err)
		assert.Equal(t, 3, len(sarees))
		assert.Equal(t, "cotton", sarees[0].FabricType)
		assert.Equal(t, "red", sarees[0].Color)
		assert.Equal(t, "saree", sarees[0].Category)
	})

	t.Run("Should save saree with a generated id and find it", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		sareeRepo := saree_repo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		saved, err := sareeService.Save(domain.Saree{
			Name:       "Red cotton",
			FabricType: "cotton",
			Color:      "red",
			Category:   "saree",
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, saved.UID)

		saree, err := sareeService.Find(saved.UID.String())
		require.NoError(t, err)
		assert.Equal(t, "Red cotton", saree.Name)
		assert.Equal(t, "cotton", saree.FabricType)
		assert.Equal(t, "red", saree.Color)
		assert.Equal(t, "saree", saree.Category)
	})

	t.Run("Should reject a saree without a name", func(t *testing.T) {

		sareeService := services.NewSareeService(saree_repo.NewSareeRepository(repo))
		_, err := sareeService.Save(domain.Saree{FabricType: "cotton"})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
	})

	t.Run("Should update saree", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		sareeRepo := saree_repo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		saved, err := sareeService.Save(domain.Saree{Name: "Saree", Color: "red"})
		require.NoError(t, err)

		updated, err := sareeService.Update(saved.UID.String(), domain.Saree{
			UID:   uuid.New(), // ignored: the id comes from the URL
			Name:  "Saree",
			Color: "blue",
			Stock: 4,
		})
		require.NoError(t, err)
		assert.Equal(t, saved.UID, updated.UID)

		saree, err := sareeService.Find(saved.UID.String())
		require.NoError(t, err)
		assert.Equal(t, "blue", saree.Color)
		assert.Equal(t, 4, saree.Stock)
	})

	t.Run("Should report a missing saree", func(t *testing.T) {

		sareeService := services.NewSareeService(saree_repo.NewSareeRepository(repo))
		missing := uuid.NewString()

		_, err := sareeService.Find(missing)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		_, err = sareeService.Update(missing, domain.Saree{Name: "x"})
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.ErrorIs(t, sareeService.Delete(missing), domain.ErrNotFound)
		_, err = sareeService.Find("not-a-uuid")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Should delete saree", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		sareeRepo := saree_repo.NewSareeRepository(repo)
		sareeService := services.NewSareeService(sareeRepo)
		saved, err := sareeService.Save(domain.Saree{Name: "Saree"})
		require.NoError(t, err)

		require.NoError(t, sareeService.Delete(saved.UID.String()))
		_, err = sareeService.Find(saved.UID.String())
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Should backfill ids for sarees saved without one", func(t *testing.T) {

		err := repo.Drop(context.Background())
		require.NoError(t, err)

		_, err = repo.InsertMany(context.Background(), []interface{}{
			bson.D{{Key: "name", Value: "no uid field"}},
			bson.D{{Key: "name", Value: "zero uid"}, {Key: "uid", Value: uuid.Nil}},
		})
		require.NoError(t, err)

		sareeRepo := saree_repo.NewSareeRepository(repo)
		fixed, err := sareeRepo.BackfillUIDs(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 2, fixed)

		sarees, err := sareeRepo.FindAll()
		require.NoError(t, err)
		require.Len(t, sarees, 2)
		assert.NotEqual(t, uuid.Nil, sarees[0].UID)
		assert.NotEqual(t, sarees[0].UID, sarees[1].UID)
	})

}
