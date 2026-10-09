package uploadhdl

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Smallest valid PNG (1×1 transparent pixel).
var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

func newRouter(t *testing.T, publicBaseURL string) (*gin.Engine, string) {
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	h, err := NewUploadHandler(dir, publicBaseURL)
	require.NoError(t, err)
	router := gin.New()
	router.POST("/uploads", h.Upload)
	router.GET("/uploads/:name", h.Serve)
	return router, dir
}

func upload(router *gin.Engine, filename string, content []byte) *httptest.ResponseRecorder {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	part, _ := mw.CreateFormFile("file", filename)
	part.Write(content)
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/uploads", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Host = "api.local:8080"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestUploadStoresImageAndServesIt(t *testing.T) {
	router, dir := newRouter(t, "")
	w := upload(router, "photo.png", tinyPNG)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var resp struct{ URL string }
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, strings.HasPrefix(resp.URL, "http://api.local:8080/uploads/"), resp.URL)
	assert.True(t, strings.HasSuffix(resp.URL, ".png"))

	name := filepath.Base(resp.URL)
	stored, err := os.ReadFile(filepath.Join(dir, name))
	require.NoError(t, err)
	assert.Equal(t, tinyPNG, stored)

	get := httptest.NewRecorder()
	router.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/uploads/"+name, nil))
	assert.Equal(t, http.StatusOK, get.Code)
	assert.Equal(t, "nosniff", get.Header().Get("X-Content-Type-Options"))
}

func TestUploadUsesPublicBaseURL(t *testing.T) {
	router, _ := newRouter(t, "https://api.example.com/")
	w := upload(router, "photo.png", tinyPNG)
	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"https://api.example.com/uploads/`)
}

func TestUploadRejectsNonImages(t *testing.T) {
	router, dir := newRouter(t, "")
	// The extension lies; the content is sniffed.
	w := upload(router, "evil.png", []byte("<html><script>alert(1)</script></html>"))
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)

	entries, _ := os.ReadDir(dir)
	assert.Empty(t, entries)
}

func TestServeOnlyImageFiles(t *testing.T) {
	router, dir := newRouter(t, "")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "notes.html"), []byte("<p>hi</p>"), 0o644))

	for _, path := range []string{"/uploads/notes.html", "/uploads/missing.png", "/uploads/..%2fsecret.png"} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, http.StatusNotFound, w.Code, path)
	}
}
