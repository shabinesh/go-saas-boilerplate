package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h handlers) ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT token from cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key used for signing
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(401)
			return
		}

		// Check token expiration
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			exp, ok := claims["exp"].(float64)
			if !ok || float64(time.Now().Unix()) > exp {
				c.AbortWithStatus(401)
				return
			}
			// Add user info to context
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
