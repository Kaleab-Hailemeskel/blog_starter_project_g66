package controllers

import (
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogController struct {
	BlogUseCase domain.IBlogUseCase
	UserUseCase domain.IUserUseCase 
}

var queryStrings = []string{"tag", "author", "title", "popularity", "date", "p"}

func NewController(blogUseCase domain.IBlogUseCase, userUseCase domain.IUserUseCase) *BlogController {
	return &BlogController{
		BlogUseCase: blogUseCase,
		UserUseCase: userUseCase,
	}
}

func (cntrl *BlogController) CreateBlog(ctx *gin.Context) {
	var blogDTO domain.BlogDTO
	if err := ctx.ShouldBindJSON(&blogDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid blog data", "details": err.Error()})
		return
	}
	blog := conv.ChangeToDomainBlog(&blogDTO)
	emailVal, exists := ctx.Get("email")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User email not found in context"})
		return
	}
	ownerEmail, ok := emailVal.(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Email in context is not a string", "value": fmt.Sprintf("%v", emailVal)})
		return
	}
	
	createdBlog, err := cntrl.BlogUseCase.CreateBlog(blog, ownerEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog", "details": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "blog created", "blog": createdBlog})
}

// filter blog can also be used to get all blogs.
func (cntrl *BlogController) FilterBlog(ctx *gin.Context) {
	mapQuery := map[string]string{}
	for _, val := range queryStrings {
		mapQuery[val] = ctx.Query(val)
	}
	//! Populartiy_value and date filter aren't implemented
	filter := domain.Filter{
		Tag:        mapQuery[queryStrings[0]],
		Title:      mapQuery[queryStrings[2]],
		AuthorName: mapQuery[queryStrings[1]],
	}
	res_int, res_err := strconv.Atoi(mapQuery[queryStrings[3]])
	if res_err == nil {
		filter.Popularity_value = res_int
	}
	res, err := strconv.Atoi(mapQuery["p"])
	if err != nil {
		res = 1
	}
	resList, err := cntrl.BlogUseCase.GetAllBlogsByFilter(&filter, res)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "blog list not found"})
		return
	}
	//? DataChange from []*BlogDTO to []*Blog
	domainBlogList := []domain.Blog{}
	for _, val := range resList {
		domainBlogList = append(domainBlogList, *conv.ChangeToDomainBlog(val))
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"result": domainBlogList})

}

// ! user should be authorized with the modifcation before doing anyting
func (cntrl *BlogController) DeleteBlog(ctx *gin.Context) {
	blogIDString := ctx.Param("id")
	emailVal, exists := ctx.Get("email")

	blogIDStringObj, err := primitive.ObjectIDFromHex(blogIDString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": "blog ID isn't correct"})
	}
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized User"})
		return
	}
	userDTO, _ := cntrl.UserUseCase.GetUserByEmail(emailVal.(string))
	if userDTO == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized User"})
		return
	}
	blog, _ := cntrl.BlogUseCase.GetBlogByID(blogIDStringObj)
	if blog.OwnerID != userDTO.UserID {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized User"})
		return
	}
	err = cntrl.BlogUseCase.DeleteBlogByID(blogIDString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Blog Deleted"})
}

//! user should be authorized with the modifcation before doing anyting

func (cntrl *BlogController) UpdateBlog(ctx *gin.Context) {
	blogStringID := ctx.Param("id")
	emailVal, exists := ctx.Get("email")

	blogIDStringObj, err := primitive.ObjectIDFromHex(blogStringID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": "blog ID isn't correct"})
	}
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized User"})
		return
	}
	userDTO, _ := cntrl.UserUseCase.GetUserByEmail(emailVal.(string))
	if userDTO == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized User"})
		return
	}
	blog, _ := cntrl.BlogUseCase.GetBlogByID(blogIDStringObj)
	if blog.OwnerID != userDTO.UserID {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized User"})
		return
	}
	var blogDTO domain.BlogDTO
	if err := ctx.ShouldBindJSON(&blogDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid Blog type"})
		return
	}
	err = cntrl.BlogUseCase.UpdateBlogByID(blogStringID, conv.ChangeToDomainBlog(&blogDTO))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "blog Updated"})
}

// ? From Popularity
func (cntrl *BlogController) LikeBlog(ctx *gin.Context) {
	// Get blog ID from param
	blogIDStr := ctx.Param("blog_id")

	// Change it to primitive.ObjectID
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		return
	}

	// Get the user email from the context (assuming it's set by middleware)
	userEmail, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found"})
		return
	}
	// Send the blogID and user email to usecase.LikeBlog()
	if err := cntrl.BlogUseCase.LikeBlog(blogID, userEmail.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like blog" + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}
func (cntrl *BlogController) DisLikeBlog(ctx *gin.Context) {
	// Get blog ID from param
	blogIDStr := ctx.Param("blog_id")

	// Change it to primitive.ObjectID
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		return
	}

	// Get the user email from the context (assuming it's set by middleware)
	userEmail, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found"})
		return
	}
	// Send the blogID and user email to usecase.LikeBlog()
	if err := cntrl.BlogUseCase.DisLikeBlog(blogID, userEmail.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike blog" + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}
func (cntrl *BlogController) CommentBlog(ctx *gin.Context) {
	// Get the comment from the Json Body
	var commentDTO domain.CommentDTO
	if err := ctx.ShouldBindJSON(&commentDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "invalid comment format"})
	}
	// Get blog ID from param
	blogIDStr := ctx.Param("blog_id")
	// Change it to primitive.ObjectID
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		return
	}

	// Get the user email from the context (assuming it's set by middleware)
	userEmail, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found"})
		return
	}
	// Send the blogID and user email to usecase.LikeBlog()
	if err := cntrl.BlogUseCase.CommentBlog(userEmail.(string), &commentDTO, blogID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like blog" + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Commented successfully"})
}
