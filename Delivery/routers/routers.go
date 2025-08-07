package routers

import (
	"blog_starter_project_g66/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func Router(uc *controllers.UserController, pc *controllers.PasswordController, bc *controllers.BlogController) {
	router := gin.Default()

	router.POST("/login",uc.HandleLogin)
	router.POST("/refresh",uc.HandleRefresh)
	router.POST("/registration", uc.Registration)
	router.POST("/registration/verification",uc.RegistrationValidation )
	router.POST("/forgot_password",pc.ForgotPassword)
	router.PUT("/reset_password", pc.ResetPassword)
	router.POST("/promote_user", uc.PromoteUser)
	router.POST("/demote_user", uc.DemoteUser)
	router.POST("/logout",)
	router.PUT("/editprofile")
	blogRoutes := router.Group("/blog")
	{
		blogRoutes.POST("", bc.CreateBlog)     
		blogRoutes.GET("", bc.FilterBlog)         
		blogRoutes.GET("/filter", bc.FilterBlog)   
		blogRoutes.PUT("/:id", bc.UpdateBlog)      
		blogRoutes.DELETE("/:id", bc.DeleteBlog)   
	}
	// router.POST("/blog/sreach",)
	// router.POST("/ai",)
	// router.POST("/ai/:id",)

	router.Run()
}