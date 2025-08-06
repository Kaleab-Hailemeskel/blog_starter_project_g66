package controllers

import (
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	usecases "blog_starter_project_g66/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase *usecases.UserUsecase
}

type PromoteDemoteRequest struct {
	TargetEmail string `json:"target_email" binding:"required,email"`
}

func NewUserUsecase(uuc *usecases.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: uuc,
	}
}
func (uc *UserController) Registration(ctx *gin.Context) {

	var user *domain.UserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	err := uc.UserUsecase.HandleRegistration(conv.ChangeToDomainUser(user))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Please enter your otp to successfully register",
	})
}

func (uc *UserController) RegistrationValidation(ctx *gin.Context) {
	var user *domain.UserUnverifiedDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	if user.Email == "" || user.OTP == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	userone := conv.ChangeToDomainVerification(user)
	valid, err := uc.UserUsecase.VerifyOTP(userone.Email, userone.OTP)
	if err != nil {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"message": "Error while verifying OTP",
	})
		return
	}

	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid or expired OTP",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

func (uc *UserController) PromoteUser(ctx *gin.Context) {
	var req PromoteDemoteRequest
	 if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    actingEmail, exists := ctx.Get("user_email")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    err := uc.UserUsecase.PromoteUser(actingEmail.(string), req.TargetEmail)
    if err != nil {
        ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "User promoted to ADMIN successfully"})
}

func (uc *UserController) DemoteUser(ctx *gin.Context) {
	var req PromoteDemoteRequest
	 if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    actingEmail, exists := ctx.Get("user_email")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    err := uc.UserUsecase.PromoteUser(actingEmail.(string), req.TargetEmail)
    if err != nil {
        ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Admin demoted to user successfully"})
}