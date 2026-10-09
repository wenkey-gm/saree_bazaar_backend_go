package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCORSRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS("https://munichandrasarees.com, http://localhost:3001"))
	router.GET("/sarees", func(c *gin.Context) { c.JSON(http.StatusOK, []string{}) })
	return router
}

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sarees", nil)
	req.Header.Set("Origin", "http://localhost:3001")
	newCORSRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3001", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSIgnoresUnknownOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sarees", nil)
	req.Header.Set("Origin", "https://evil.example")
	newCORSRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSAnswersPreflight(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/sarees", nil)
	req.Header.Set("Origin", "https://munichandrasarees.com")
	newCORSRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://munichandrasarees.com", w.Header().Get("Access-Control-Allow-Origin"))
}
