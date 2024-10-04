package controllers

import (
	"log"
	"login_page/database"
	"login_page/models"
	"login_page/utils"
	"net/http"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserSignup(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "user_signup.html", gin.H{"title": "Sign Up"})
		return
	}

	if c.Request.Method == http.MethodPost {
		var Form struct {
			FirstName       string `form:"first_name" binding:"required"`
			LastName        string `form:"last_name" binding:"required"`
			Email           string `form:"email" binding:"required,email"`
			Password        string `form:"password" binding:"required"`
			ConfirmPassword string `form:"confirm_password" binding:"required"`
			PhoneNumber     string `form:"phone_number" binding:"required"`
		}

		if err := c.ShouldBind(&Form); err != nil {
			log.Printf("Form binding error: %v", err)
			c.HTML(http.StatusBadRequest, "user_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Form binding error: " + err.Error(),
			})
			return
		}

		if Form.Password != Form.ConfirmPassword {
			c.HTML(http.StatusBadRequest, "user_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Passwords do not match",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Form.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			c.HTML(http.StatusInternalServerError, "user_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Password hashing error: " + err.Error(),
			})
			return
		}

		// Generate OTP
		// otp, err := utils.GenerateOTP()
		// if err != nil {
		// 	log.Printf("Error generating OTP: %v", err)
		// 	c.HTML(http.StatusInternalServerError, "user_signup.html", gin.H{
		// 		"title": "Sign Up",
		// 		"error": "Error generating OTP: " + err.Error(),
		// 	})
		// 	return
		// }
		otp, err := utils.GenerateOTP()
		if err != nil {
			log.Printf("Error generating OTP: %v", err)
			return
		}
		log.Printf("Generated OTP: %s", otp) // Add this line

		// Set expiration time for OTP
		otpExpiration := time.Now().Add(10 * time.Minute)

		// Create the user object
		user := models.User{
			FirstName:    Form.FirstName,
			LastName:     Form.LastName,
			Email:        Form.Email,
			Password:     string(hashedPassword),
			PhoneNumber:  Form.PhoneNumber,
			Status:       "pending",
			OTPCode:      otp,
			OTPExpiresAt: otpExpiration,
		}

		// Save user to the database
		log.Printf("OTP: %s, ExpiresAt: %v", user.OTPCode, user.OTPExpiresAt)
		if err := database.CreateUser(database.DB, &user); err != nil {
			log.Printf("Database error: %v", err)
			c.HTML(http.StatusInternalServerError, "user_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Database error: " + err.Error(),
			})
			return
		}

		// Send the OTP to the user's email
		if err := utils.SendOTPEmail(user.Email, otp); err != nil {
			log.Printf("Error sending OTP email to %s: %v", user.Email, err)
			c.HTML(http.StatusInternalServerError, "user_signup.html", gin.H{
				"title": "Sign Up",
				"error": "Error sending OTP",
			})
			return
		}

		// Redirect to OTP verification page
		c.Redirect(http.StatusSeeOther, "/verify-otp?email="+user.Email)
	}
}

func EmailVerifiedMiddleware(s, verificationToken string) {
	panic("unimplemented")
}

func Login(c *gin.Context) {
	var form struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		log.Printf("Form binding error: %v", err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByEmail(form.Email)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	if user == nil {
		log.Printf("User not found: %s", form.Email)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
		log.Printf("Password mismatch for user: %s", form.Email)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	if user.Status != "active" {
		log.Printf("User status not active: %s", user.Status)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Email not verified"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: form.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error generating token"})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
	c.Set("userEmail", form.Email)

	c.Redirect(http.StatusSeeOther, "/home")
}
