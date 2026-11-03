package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"product_api/internal/core/domain"
	"product_api/internal/core/services"
	"strings"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s *services.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				c.JSON(500, gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			// otherwise error type is unknown

			c.JSON(500, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {

			c.JSON(401, gin.H{
				"error": "Authorization header must be of the form 'Bearer token'",
			})
			c.Abort()
			return
		}

		// validate ID token here
		user, err := s.ValidateIDToken(idTokenHeader[1])

		if err != nil {

			c.JSON(401, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

// RequireAdmin must run after AuthUser; it rejects users without the admin role.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")
		if u, isUser := user.(*domain.User); !ok || !isUser || u.Role != domain.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
