// Package adminui serves the catalog admin panel: a small static app that
// signs in through /login and manages sarees through the JSON API.
package adminui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFiles embed.FS

const contentSecurityPolicy = "default-src 'self'; img-src 'self' https: data: blob:; " +
	"script-src 'self'; style-src 'self'; connect-src 'self'; " +
	"frame-ancestors 'none'; base-uri 'none'; form-action 'self'; object-src 'none'"

func securityHeaders(c *gin.Context) {
	c.Header("Content-Security-Policy", contentSecurityPolicy)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("Cache-Control", "no-cache")
}

// Register mounts the admin panel at /admin.
func Register(router *gin.Engine) {
	static, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	group := router.Group("/admin", securityHeaders)
	group.GET("", func(c *gin.Context) {
		c.FileFromFS("/", http.FS(static))
	})
	group.StaticFS("/assets", http.FS(static))
}
