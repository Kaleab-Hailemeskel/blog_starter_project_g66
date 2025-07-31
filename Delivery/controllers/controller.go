package controllers

import (
	
	"net/http"
	"blog_starter_project_g66/Usecases"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase *usecases.UserUsecase
}

func NewUserUsecase(uuc *usecases.UserUsecase) *UserController{
	return &UserController{
		UserUsecase: uuc,
	}
}
func (uc *UserController) Registration(ctx *gin.Context){

	var user *UserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusNotFound,gin.H{
			"error": "Invalid request payload",
		} )
		return
	}
	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	err := uc.UserUsecase.HandleRegistration()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "successuly created",
	})
}