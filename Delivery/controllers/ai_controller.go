package controllers

import (
	domain "blog_starter_project_g66/Domain"
	infrastructure "blog_starter_project_g66/Infrastructure"
	"encoding/json"
	"log"

	// usecases "blog_starter_project_g66/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIController struct {
	AiComment domain.IAICommentUsecase
	Aiblog    domain.IAIBlogUsecase
	AiFilter  domain.IAIFilterUsecase
}

func NewAIController(aiComment domain.IAICommentUsecase, aiblog domain.IAIBlogUsecase, aiFilter domain.IAIFilterUsecase) *AIController {
	return &AIController{
		AiComment: aiComment,
		Aiblog:    aiblog,
		AiFilter:  aiFilter,
	}

}

func (ac *AIController) HandleAIComment(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	idStr, _ := userId.(string)
	objID, ok := primitive.ObjectIDFromHex(idStr)
	log.Println("start")

	if ok != nil {
		log.Printf("👀 This is the one")
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "not valid id", "id": idStr})
		return
	}
	log.Println("00000000000-----------")
	aIInteraction := infrastructure.NewAICommentInteraction(objID)
	log.Println("+++++++++", aIInteraction, "+++++++++")
	var userRequest *domain.AICommentDTO
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	// ty := aIInteraction.AICommentUsecase(userRequest,aIInteraction)
	comment, ok := ac.AiComment.AICommentUsecase(userRequest, aIInteraction)
	log.Println("1", ok, comment)
	if ok != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": ok.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": comment,
	})
}
func (ac *AIController) HandleAIBogModFromUser(ctx *gin.Context) {
	log.Println("=>> correct Way")
	userId, _ := ctx.Get("user_id")
	idStr, _ := userId.(string)
	objID, ok := primitive.ObjectIDFromHex(idStr)
	if ok != nil {
		return
	}
	var userRequest *domain.AIUserBlogDTO
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	aIInteraction := infrastructure.NewAIBlogInteraction(objID)
	str, err := json.Marshal(userRequest.BlogDTO)
	if err != nil {
		log.Println("Weel a;lkdjf;alk")
	}
	resDomain, err := aIInteraction.CallAIAndGetResponse("", userRequest.UserMessage, string(str))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "AI can't return message"})
	}
	finalRes := aIInteraction.ParseJsonBodyToDomain(resDomain)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message":finalRes})
}
func (ac *AIController) HandleAIBog(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	idStr, _ := userId.(string)
	objID, ok := primitive.ObjectIDFromHex(idStr)
	if ok != nil {
		return
	}
	aIInteraction := infrastructure.NewAIBlogInteraction(objID)

	blogId := ctx.Param("id")
	var userRequest *domain.AIBlogDTO
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	response, err := ac.Aiblog.AIBlogUsecase(blogId, userRequest, aIInteraction)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": response,
	})
}

func (ac *AIController) HandleAIFilter(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	idStr, _ := userId.(string)
	objID, ok := primitive.ObjectIDFromHex(idStr)
	if ok != nil {
		return
	}
	aIInteraction := infrastructure.NewAIBlogFilterInteraction(objID)
	var userRequest *domain.AIBlogDTO
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	response, err := ac.AiFilter.AIFilterUsecase(userRequest, aIInteraction)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": response,
	})
}
