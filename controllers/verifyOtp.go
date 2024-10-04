package controllers

import (
	"login_page/database"
	"login_page/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func VerifyOTP(c *gin.Context) {
// 	if c.Request.Method == http.MethodGet {
// 		email := c.Query("email")
// 		c.HTML(http.StatusOK, "verify_otp.html", gin.H{
// 			"email": email,
// 		})
// 		return
// 	}

// 	if c.Request.Method == http.MethodPost {
// 		var Form struct {
// 			OTP   string `form:"otp" binding:"required"`
// 			Email string `form:"email" binding:"required"`
// 		}

// 		if err := c.ShouldBind(&Form); err != nil {
// 			c.HTML(http.StatusBadRequest, "verify_otp.html", gin.H{
// 				"email": Form.Email,
// 				"error": "Form binding error: " + err.Error(),
// 			})
// 			return
// 		}

// 		user, err := database.GetUserByEmail(Form.Email)
// 		if err != nil {
// 			c.HTML(http.StatusBadRequest, "verify_otp.html", gin.H{
// 				"email": Form.Email,
// 				"error": "User not found",
// 			})
// 			return
// 		}

// 		if user.OTPCode != Form.OTP || time.Now().After(user.OTPExpiresAt) {
// 			c.HTML(http.StatusBadRequest, "verify_otp.html", gin.H{
// 				"email": Form.Email,
// 				"error": "Invalid or expired OTP",
// 			})
// 			return
// 		}

// 		// OTP is valid; update the user's status to "active"
// 		user.Status = "active"
// 		if err := database.UpdateUser(database.DB, user); err != nil {
// 			c.HTML(http.StatusInternalServerError, "verify_otp.html", gin.H{
// 				"email": Form.Email,
// 				"error": "Error updating user status",
// 			})
// 			return
// 		}

// 		// Generate JWT for the user
// 		token, err := middleware.GenerateJWT(*user)
// 		if err != nil {
// 			c.HTML(http.StatusInternalServerError, "verify_otp.html", gin.H{
// 				"email": Form.Email,
// 				"error": "Error generating token",
// 			})
// 			return
// 		}

// 		// Set JWT in a cookie or send it back in the response
// 		c.SetCookie("token", token, 3600, "/", "", false, true)

//			// Redirect to the user's dashboard
//			c.Redirect(http.StatusSeeOther, "/home")
//		}
//	}
func VerifyOTP(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		email := c.Query("email")
		c.HTML(http.StatusOK, "verify_otp.html", gin.H{
			"email": email,
		})
		return
	}

	if c.Request.Method == http.MethodPost {
		var Form struct {
			OTP   string `form:"otp" binding:"required"`
			Email string `form:"email"`
		}

		if err := c.ShouldBind(&Form); err != nil {
			c.HTML(http.StatusBadRequest, "verify_otp.html", gin.H{
				"email": Form.Email,
				"error": "Form binding error: " + err.Error(),
			})
			return
		}

		// Fetch the user by email
		user, err := database.GetUserByEmail(Form.Email)
		if err != nil || user == nil {
			c.HTML(http.StatusBadRequest, "verify_otp.html", gin.H{
				"email": Form.Email,
				"error": "User not found",
			})
			return
		}

		// Check OTP validity
		if user.OTPCode != Form.OTP || time.Now().After(user.OTPExpiresAt) {
			c.HTML(http.StatusBadRequest, "verify_otp.html", gin.H{
				"email": Form.Email,
				"error": "Invalid or expired OTP",
			})
			return
		}

		// OTP is valid, update the user's status to active
		user.Status = "active"
		if err := database.UpdateUser(database.DB, user); err != nil {
			c.HTML(http.StatusInternalServerError, "verify_otp.html", gin.H{
				"email": Form.Email,
				"error": "Error updating user status",
			})
			return
		}

		// Generate JWT for the user
		token, err := middleware.GenerateJWT(*user)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "verify_otp.html", gin.H{
				"email": Form.Email,
				"error": "Error generating token",
			})
			return
		}

		// Set JWT in a cookie or send it back in the response
		c.SetCookie("token", token, 3600, "/", "", false, true)

		// Redirect to the user's dashboard
		c.Redirect(http.StatusSeeOther, "/home")
	}
}
