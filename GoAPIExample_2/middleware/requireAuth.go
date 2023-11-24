package middleware

import (
	"Module/API/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {

	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")

	// The header should be in the format "Bearer <token>"
	// So we split the header into two parts
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// The second part of the header is the token
	tokenString := parts[1]

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("qwertyuiopasdfghjklzxcvbnm123456"), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract the user ID from the subject field in the token header
	userID := token.Header["sub"]

	// Find the user with the extracted user ID
	var user models.User
	models.DB.First(&user, userID)

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Attach the user to the request context
	c.Set("user", user)

	// Continue with the request
	c.Next()

}
