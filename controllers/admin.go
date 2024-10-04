package controllers

import (
	"login_page/database"
	"login_page/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AdminSignup(c *gin.Context) {

	var Form struct {
		FirstName       string `form:"first_name" binding:"required"`
		LastName        string `form:"last_name" binding:"required"`
		Email           string `form:"email" binding:"required,email"`
		Password        string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required"`
		PhoneNumber     string `form:"phone_number" binding:"required"`
	}

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "admin_signup.html", gin.H{"title": "Sign Up"})
	} else if c.Request.Method == http.MethodPost {

		// Bind form data
		if err := c.ShouldBind(&Form); err != nil {
			c.HTML(http.StatusBadRequest, "admin_signup.html", gin.H{
				"title": "Sign Up",
				"error": err.Error(),
			})
			return
		}

		// Validate passwords match
		if Form.Password != Form.ConfirmPassword {
			c.HTML(http.StatusBadRequest, "admin_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Passwords do not match",
			})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Form.Password), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "admin_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Failed to hash password",
			})
			return
		}

		// Map Form data to the User model
		admin := models.User{
			FirstName:   Form.FirstName,
			LastName:    Form.LastName,
			Email:       Form.Email,
			Password:    string(hashedPassword), // Save hashed password
			PhoneNumber: Form.PhoneNumber,
		}

		// Set a default status if it's not provided
		if admin.Status == "" {
			admin.Status = "active" // or whatever the default should be
		}

		// Call the CreateUser function from the 'database' package
		if err := database.CreateAdmin(database.DB, &admin); err != nil {
			c.HTML(http.StatusInternalServerError, "user_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Failed to create user",
			})
			return
		}

		// Redirect to login page after successful signup
		c.Redirect(http.StatusSeeOther, "/adminlogin")
	}

}
func AdminLogin(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		// Render the login form
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
			"title": "Login",
		})
		return
	}

	if c.Request.Method == http.MethodPost {
		// Handle form submission
		var form struct {
			Email    string `form:"email" binding:"required,email"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.ShouldBind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "adminlogin.html", gin.H{
				"title": "Login",
				"error": err.Error(),
			})
			return
		}

		// Fetch the user from the database
		admin, err := database.GetAdminByEmail(form.Email)
		if err != nil || admin == nil {
			c.HTML(http.StatusUnauthorized, "addmin_login.html", gin.H{
				"title": "Login",
				"error": "Invalid email or password.",
			})
			return
		}

		// Compare the provided password with the stored hashed password
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(form.Password))
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"title": "Login",
				"error": "Invalid email or password.",
			})
			return
		}

		// Generate JWT or set session and redirect to user home or dashboard
		// ...

		c.Redirect(http.StatusSeeOther, "/home")
	}
}
