package middleware

import (
	"log"
	"login_page/database"
	"login_page/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func EmailVerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user email from the context
		email, exists := c.Get("userEmail")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Fetch the user from the database
		user, err := database.GetUserByEmail(email.(string))
		if err != nil {
			log.Printf("Database error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Check if the user's status is "active" (i.e., email is verified via OTP)
		if user.Status != "active" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email not verified. Please verify your OTP."})
			return
		}

		// Proceed with the next middleware/handler
		c.Next()
	}
}

var jwtSecret = []byte("your_secret_key")

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
