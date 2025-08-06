package routers

import (
	"blog_starter_project_g66/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func Router(uc *controllers.UserController, pc *controllers.PasswordController, ctl *controllers.UserController) {
	router := gin.Default()

	router.POST("/login",uc.HandleLogin)
	router.POST("/refresh",uc.HandleRefresh)
	router.POST("/registration", uc.Registration)
	router.POST("/registration/verification",uc.RegistrationValidation )
	router.POST("/forgot_password",pc.ForgotPassword)
	router.PUT("/reset_password", pc.ResetPassword)
	router.POST("/promote_user", ctl.PromoteUser)
	router.POST("/demote_user", ctl.DemoteUser)
	router.POST("/logout",)
	router.PUT("/editprofile")
	router.POST("/blog",)
	router.GET("/blog",)
	router.GET("/blog/filter",)
	// router.POST("/foget_password",)
	router.POST("/logout",uc.HandleLogout)
	// router.PUT("/editprofile")
	// router.POST("/blog",)
	// router.GET("/blog",)
	// router.GET("/blog/filter",)
	// router.PUT("/blog/:id",)
	// router.DELETE("/blog/:id",)
	// router.POST("/blog/sreach",)
	// router.POST("/ai",)
	// router.POST("/ai/:id",)

	router.Run()
}