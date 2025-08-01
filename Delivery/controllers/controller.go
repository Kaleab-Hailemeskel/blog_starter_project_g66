package controllers

import (
	domain "blog_starter_project_g66/Domain"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	BlogUseCase domain.IBlogUseCase
	// UserUseCase domain.IUserUseCase  //! For the time being it is commented out
}

func NewController(blogUseCase domain.IBlogUseCase) *Controller {
	return &Controller{
		BlogUseCase: blogUseCase,
	}
}

func (cntrl *Controller) FilterBlog(ctx *gin.Context) {
	ctx.Query("")
}
