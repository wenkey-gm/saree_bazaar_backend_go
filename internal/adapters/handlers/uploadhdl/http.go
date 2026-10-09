package uploadhdl

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MaxImageBytes is the largest photo accepted (phone photos are ~3–8 MB).
const MaxImageBytes = 15 << 20

// Detected content type → stored file extension. Only formats every browser
// renders are accepted; the type is sniffed from the bytes, not the filename.
var allowedTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type UploadHandler struct {
	dir           string
	publicBaseURL string
}

// NewUploadHandler stores images in dir. Returned URLs start with
// publicBaseURL (e.g. "https://api.example.com"); when it's empty they're
// built from the request's host.
func NewUploadHandler(dir, publicBaseURL string) (*UploadHandler, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &UploadHandler{dir: dir, publicBaseURL: strings.TrimRight(publicBaseURL, "/")}, nil
}

func (h *UploadHandler) baseURL(c *gin.Context) string {
	if h.publicBaseURL != "" {
		return h.publicBaseURL
	}
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// Upload accepts a multipart "file" field and responds with {"url": ...}.
func (h *UploadHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxImageBytes+1<<20)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "send the image in a \"file\" form field (max 15 MB)"})
		return
	}
	defer file.Close()
	if header.Size > MaxImageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "image is larger than 15 MB"})
		return
	}

	head := make([]byte, 512)
	n, err := io.ReadFull(file, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read image"})
		return
	}
	ext, ok := allowedTypes[http.DetectContentType(head[:n])]
	if !ok {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "only JPEG, PNG, WebP or GIF images are allowed"})
		return
	}

	name := uuid.NewString() + ext
	out, err := os.OpenFile(filepath.Join(h.dir, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		log.Printf("could not create upload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save image"})
		return
	}
	_, err = io.Copy(out, io.MultiReader(bytes.NewReader(head[:n]), file))
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(filepath.Join(h.dir, name))
		log.Printf("could not write upload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save image"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url": h.baseURL(c) + "/uploads/" + name})
}

// Serve returns a handler for GET /uploads/:name. Directory listings are not
// exposed, and files are served with nosniff so they're only ever images.
func (h *UploadHandler) Serve(c *gin.Context) {
	name := filepath.Base(c.Param("name"))
	if _, ok := allowedTypes[mimeForExt(filepath.Ext(name))]; !ok {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.File(filepath.Join(h.dir, name))
}

func mimeForExt(ext string) string {
	for mime, e := range allowedTypes {
		if e == ext {
			return mime
		}
	}
	return ""
}
