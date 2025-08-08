package routers

import (
	"blog_starter_project_g66/Delivery/controllers"
	infrastructure "blog_starter_project_g66/Infrastructure"

	"github.com/gin-gonic/gin"
)

func Router(uc *controllers.UserController, pc *controllers.PasswordController, bc *controllers.BlogController, auth *infrastructure.AuthMiddleware) {
	router := gin.Default()

	router.POST("/login",uc.HandleLogin)
	router.POST("/refresh",uc.HandleRefresh)
	router.POST("/registration", uc.Registration)
	router.POST("/registration/verification",uc.RegistrationValidation )
	router.POST("/forgot_password",pc.ForgotPassword)
	router.PUT("/reset_password", pc.ResetPassword)

	blogRoutes := router.Group("/blog")
	blogRoutes.Use(auth.JWTAuthMiddleware())
	{
		blogRoutes.POST("", bc.CreateBlog)     
		blogRoutes.GET("", bc.FilterBlog)         
		blogRoutes.GET("/filter", bc.FilterBlog)   
		blogRoutes.PUT("/:id", bc.UpdateBlog)      
		blogRoutes.DELETE("/:id", bc.DeleteBlog)   
	}

	adminRoutes := router.Group("/")
	adminRoutes.Use(auth.JWTAuthMiddleware(), infrastructure.RoleMiddleware("SUPER_ADMIN"))
	{
		adminRoutes.POST("/promote_user", uc.PromoteUser)
		adminRoutes.POST("/demote_user", uc.DemoteUser)
	}

	userRoutes := router.Group("/user")
	userRoutes.Use(auth.JWTAuthMiddleware())
	{
		userRoutes.POST("/logout",uc.HandleLogout)
		userRoutes.PUT("/profile", uc.UpdateProfile)
	}
	// router.POST("/blog/sreach",)
	// router.POST("/ai",)
	// router.POST("/ai/:id",)

	router.Run()
}