package controllers

import (
	"blog_starter_project_g66/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordController struct {
	passwordUc domain.IPasswordUsecase
}

func NewPasswordController(uc domain.IPasswordUsecase) *PasswordController {
	return &PasswordController{passwordUc: uc}
}

func (pc *PasswordController) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := pc.passwordUc.GenerateResetToken(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (pc *PasswordController) ResetPassword(c *gin.Context) {
	var req struct {
		Token 		string `json:"token"`
		NewPassword	string	`json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := pc.passwordUc.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
} 