package sareehdl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func GetTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	body := `{"fabric_type":"butta","category":"saree","color":"blue"}`
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		Body:   io.NopCloser(strings.NewReader(body)),
	}

	return ctx
}

func MockJson(c *gin.Context, params gin.Params, u url.Values, method string) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")

	// set path params
	c.Params = params
}

type MockSareeService struct {
}

// Delete implements ports.SareeService
func (*MockSareeService) Delete(id string) error {
	return nil
}

func (m *MockSareeService) FindAll() ([]domain.Saree, error) {
	return []domain.Saree{}, nil
}

func (m *MockSareeService) Find(id string) (domain.Saree, error) {
	return domain.Saree{}, nil
}

func (m *MockSareeService) Save(saree domain.Saree) (domain.Saree, error) {
	return domain.Saree{}, nil
}

func (m *MockSareeService) Update(id string, saree domain.Saree) (domain.Saree, error) {
	return domain.Saree{}, nil
}

func NewMockSareeService() *MockSareeService {
	return &MockSareeService{}
}

func NewMockSareeHandler() *SareeHandler {
	return NewSareeHandler(NewMockSareeService())
}

func TestSareeHandler_FindAll(t *testing.T) {
	type fields struct {
		sareeService ports.SareeService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestSareeHandler_FindAll",
			fields: fields{
				sareeService: NewMockSareeService(),
			},
			args: args{
				c: GetTestGinContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SareeHandler{
				sareeService: tt.fields.sareeService,
			}
			s.FindAll(tt.args.c)

			assert.Equal(t, 200, tt.args.c.Writer.Status())

			var sarees []domain.Saree
			tt.args.c.BindJSON(&sarees)
			assert.Equal(t, 0, len(sarees))
		})
	}
}

func TestSareeHandler_Find(t *testing.T) {
	type fields struct {
		sareeService ports.SareeService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestSareeHandler_Find",
			fields: fields{
				sareeService: NewMockSareeService(),
			},
			args: args{
				c: GetTestGinContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SareeHandler{
				sareeService: tt.fields.sareeService,
			}

			params := gin.Params{
				{
					Key:   "id",
					Value: "643a5b3f4be6ee2f5ab5a940",
				},
			}

			MockJson(tt.args.c, params, url.Values{}, "GET")

			s.Find(tt.args.c)

			assert.Equal(t, 200, tt.args.c.Writer.Status())

			var saree domain.Saree
			tt.args.c.BindJSON(&saree)
			assert.Equal(t, "butta", saree.FabricType)
		})
	}
}

func TestSareeHandler_Save(t *testing.T) {
	type fields struct {
		sareeService ports.SareeService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestSareeHandler_Save",
			fields: fields{
				sareeService: NewMockSareeService(),
			},
			args: args{
				c: GetTestGinContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SareeHandler{
				sareeService: tt.fields.sareeService,
			}
			s.Save(tt.args.c)

			assert.Equal(t, 200, tt.args.c.Writer.Status())

			var saree domain.Saree
			tt.args.c.BindJSON(&saree)
			assert.Equal(t, "", saree.FabricType)
		})
	}
}

func TestSareeHandler_Update(t *testing.T) {
	type fields struct {
		sareeService ports.SareeService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestSareeHandler_Update",
			fields: fields{
				sareeService: NewMockSareeService(),
			},
			args: args{
				c: GetTestGinContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SareeHandler{
				sareeService: tt.fields.sareeService,
			}
			s.Update(tt.args.c)

			assert.Equal(t, 200, tt.args.c.Writer.Status())

			var saree domain.Saree
			tt.args.c.BindJSON(&saree)
			assert.Equal(t, "", saree.FabricType)
		})
	}
}

func TestSareeHandler_Delete(t *testing.T) {
	type fields struct {
		sareeService ports.SareeService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestSareeHandler_Delete",
			fields: fields{
				sareeService: NewMockSareeService(),
			},
			args: args{
				c: GetTestGinContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SareeHandler{
				sareeService: tt.fields.sareeService,
			}

			params := gin.Params{
				{
					Key:   "id",
					Value: "643a5b3f4be6ee2f5ab5a940",
				},
			}

			MockJson(tt.args.c, params, url.Values{}, "DELETE")

			s.Delete(tt.args.c)

			assert.Equal(t, 200, tt.args.c.Writer.Status())
		})
	}
}
