package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	// the JWT secret is a random 32 byte string to improve security
	// since the attacker could potentially brute force the token
	secret := []byte(os.Getenv("JWT_SECRET"))

	return func(c *gin.Context) {
		// check if in the request there is an authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			return
		}

		// the authorization header must be a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header format must be Bearer {token}",
			})
			return
		}

		// extract the token from the authorization header
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				alg, _ := token.Header["alg"].(string)
				return nil, errors.New("unexpected signing method: " + alg)
    	}
			return secret, nil
		})

		// validate the token
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token: " + err.Error(),
			})
			return
		}

		// get the user ID from the token claims and set it in the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sub, ok := claims["sub"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token claims",
				})
				return
			}
			c.Set("userID", sub)
		}

		c.Next()
	}
}