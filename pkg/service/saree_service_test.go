package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saree_bazaar.com/pkg/domain/modal"
	"saree_bazaar.com/pkg/usecase/repository"
)

func TestNewSareeService(t *testing.T) {
	tests := []struct {
		name string
		want repository.SareeRepository
	}{
		{
			name: "Test NewSareeService",
			want: &sareeService{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewSareeService())
		})
	}

}

func TestSareeServiceGetAllSarees(t *testing.T) {

	tests := []struct {
		name    string
		want    []modal.Saree
		wantErr bool
	}{
		{
			name: "Test SareeServiceGetAllSarees",
			want: []modal.Saree{
				{
					Name:  "Saree 1",
					Price: 100,
					Image: byte(1),
					Type:  "Type 1",
					Color: "Color 1",
				},
				{
					Name:  "Saree 2",
					Price: 200,
					Image: byte(2),
					Type:  "Type 2",
					Color: "Color 2",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockSareeRepository{}
			got, err := s.GetAllSarees()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sareeService_GetSaree(t *testing.T) {
	type args struct {
		id primitive.ObjectID
	}
	tests := []struct {
		name    string
		s       *sareeService
		args    args
		want    modal.Saree
		wantErr bool
	}{
		{
			name: "Test sareeService.GetSaree",
			s:    &sareeService{},
			args: args{
				id: primitive.NewObjectID(),
			},
			want: modal.Saree{
				Name:  "Saree 1",
				Price: 100,
				Image: byte(1),
				Type:  "Type 1",
				Color: "Color 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockSareeRepository{}
			got, err := s.GetSaree(tt.args.id.Hex())
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

type mockSareeRepository struct {
}

func (m *mockSareeRepository) GetAllSarees() ([]modal.Saree, error) {
	return []modal.Saree{
		{
			Name:  "Saree 1",
			Price: 100,
			Image: byte(1),
			Type:  "Type 1",
			Color: "Color 1",
		},
		{
			Name:  "Saree 2",
			Price: 200,
			Image: byte(2),
			Type:  "Type 2",
			Color: "Color 2",
		},
	}, nil
}

func (m *mockSareeRepository) GetSaree(id string) (modal.Saree, error) {
	return modal.Saree{
		Name:  "Saree 1",
		Price: 100,
		Image: byte(1),
		Type:  "Type 1",
		Color: "Color 1",
	}, nil
}
