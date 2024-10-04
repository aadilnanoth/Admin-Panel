package routes

import (
	"login_page/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	// Home route
	//  r.GET("/", controllers.UserHome)

	// Signup routes
	r.GET("/signup", controllers.UserSignup)  // Display the signup form
	r.POST("/signup", controllers.UserSignup) // Handle form submission

	//  Login routes
	r.GET("/login", controllers.Login)  // Display the login form
	r.POST("/login", controllers.Login) // Handle login form submission

	r.GET("/verify-otp", controllers.VerifyOTP)
	r.POST("/verify-otp", controllers.VerifyOTP)

	// Home route
	r.GET("/home", controllers.Home)

}
