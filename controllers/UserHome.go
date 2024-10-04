package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	// Get JWT from cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Parse and validate token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Render home page
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "User Home",
		"email": claims.Email, // Pass user's email to display on the page
	})
}
