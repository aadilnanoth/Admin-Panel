package routes

import (
	"login_page/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoute(r *gin.Engine) {

	r.GET("/adminsignup", controllers.AdminSignup)
	r.POST("/adminsignup", controllers.AdminSignup)

	r.GET("/adminlogin", controllers.AdminLogin)
	r.POST("/adminlogin", controllers.AdminLogin)

	// admin := r.Group("/admin")

	// admin.GET("/", controllers.AdminHome)
	// admin.POST("/add", controllers.AddUser)
	// admin.POST("/edit/:id", controllers.EditUser)
	// admin.POST("/block/:id", controllers.BlockUser)
	// admin.POST("/delete/:id", controllers.DeleteUser)
	// admin.GET("/search", controllers.SearchUsers)
}
