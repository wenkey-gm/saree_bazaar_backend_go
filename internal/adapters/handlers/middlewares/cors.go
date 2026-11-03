package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS allows browser requests from the comma-separated origins in
// allowedOrigins (e.g. "https://munichandrasarees.com,http://localhost:3001").
// Requests from other origins get no CORS headers, so browsers block them;
// server-to-server calls are unaffected.
func CORS(allowedOrigins string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, origin := range strings.Split(allowedOrigins, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Max-Age", "600")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
