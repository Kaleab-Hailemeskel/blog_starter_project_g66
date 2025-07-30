package routers

import "github.com/gin-gonic/gin"

func Router() {
	router := gin.Default()

	router.POST("/login",)
	router.POST("/registration",)
	router.POST("/foget_password",)
	router.POST("/logout",)
	router.PUT("/editprofile")
	router.POST("/blog",)
	router.GET("/blog",)
	router.GET("/blog/filter",)
	router.PUT("/blog/:id",)
	router.DELETE("/blog/:id",)
	router.POST("/blog/sreach",)
	router.POST("/ai",)
	router.POST("/ai/:id",)

	router.Run()
}