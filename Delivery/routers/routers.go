package routers

import (
	"blog_starter_project_g66/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func Router(uc *controllers.UserController) {
	router := gin.Default()

	router.POST("/login",)
	router.POST("/registration", uc.Registration)
	router.POST("/registration/verification",uc.RegistrationValidation )
	router.POST("/foget_password",)
	router.POST("/logout",)
	router.PUT("/editprofile")
	router.POST("/blog",)
	router.GET("/blog",)
	router.GET("/blog/filter",)
	// router.PUT("/blog/:id",)
	// router.DELETE("/blog/:id",)
	// router.POST("/blog/sreach",)
	// router.POST("/ai",)
	// router.POST("/ai/:id",)

	router.Run()
}