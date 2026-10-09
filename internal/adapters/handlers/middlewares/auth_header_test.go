package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"product_api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func requireAdminStatus(user interface{}) int {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/sarees", func(c *gin.Context) {
		if user != nil {
			c.Set("user", user)
		}
		c.Next()
	}, RequireAdmin(), func(c *gin.Context) { c.Status(http.StatusCreated) })

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/sarees", nil))
	return w.Code
}

func TestRequireAdmin(t *testing.T) {
	assert.Equal(t, http.StatusCreated, requireAdminStatus(&domain.User{Role: domain.RoleAdmin}))
	assert.Equal(t, http.StatusForbidden, requireAdminStatus(&domain.User{Role: domain.RoleCustomer}))
	assert.Equal(t, http.StatusForbidden, requireAdminStatus(nil))
}
