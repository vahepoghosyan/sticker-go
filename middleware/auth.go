package middleware

import (
	"net/http"
	"sticker-go/config"
	"sticker-go/views"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, views.ResponseError("Missing token"))
			c.Abort()
			return
		}

		claims, err := config.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, views.ResponseError("Invalid token"))
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}