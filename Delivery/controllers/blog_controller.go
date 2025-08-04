package controllers

import (
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	BlogUseCase domain.IBlogUseCase
	// UserUseCase domain.IUserUseCase  //! For the time being it is commented out
}

var queryStrings = []string{"tag", "author", "title", "popularity", "date", "p"}

func NewController(blogUseCase domain.IBlogUseCase) *BlogController {
	return &BlogController{
		BlogUseCase: blogUseCase,
	}
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
	res, err := strconv.Atoi(mapQuery["p"])
	if err != nil {
		res = 1
	}
	resList, err := cntrl.BlogUseCase.GetAllBlogsByFilter(&filter, res)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "blog list not found"})
		return
	}
	//? DataChange from []*BlogDTO to Blog
	domainBlogList := []domain.Blog{}
	for _, val := range resList {
		domainBlogList = append(domainBlogList, *conv.ChangeToDomainBlog(val))
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"result": domainBlogList})

}
func (cntrl *BlogController) DeleteBlog(ctx *gin.Context) {
	// blogObjID, err := primitive.ObjectIDFromHex(blogIDString)
	blogIDString := ctx.Param("id")

	err := cntrl.BlogUseCase.DeleteBlogByID(blogIDString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Blog Deleted"})
}
func (cntrl *BlogController) UpdateBlog(ctx *gin.Context) {
	blogStringID := ctx.Param("id")
	var blogDTO domain.BlogDTO
	if err := ctx.ShouldBindBodyWithJSON(blogDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid Blog type"})
		return
	}
	err := cntrl.BlogUseCase.UpdateBlogByID(blogStringID, conv.ChangeToDomainBlog(&blogDTO))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message":"blog Updated"})
}
