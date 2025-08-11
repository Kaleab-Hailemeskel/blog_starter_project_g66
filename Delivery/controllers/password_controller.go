package controllers

import (
	domain "blog_starter_project_g66/Domain"
	infrastructure "blog_starter_project_g66/Infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordController struct {
	passwordUc domain.IPasswordUsecase
	validate   *infrastructure.PasswordService
}

func NewPasswordController(uc domain.IPasswordUsecase, validate_ *infrastructure.PasswordService) *PasswordController {
	return &PasswordController{
		passwordUc: uc,
		validate:   validate_,
	}
}

func (pc *PasswordController) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := pc.passwordUc.GenerateResetToken(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset link sent",
	})
}

func (pc *PasswordController) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !pc.validate.IsStrongPassword(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weak password choosen, please make your password stronger"})
		return
	}

	if err := pc.passwordUc.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
