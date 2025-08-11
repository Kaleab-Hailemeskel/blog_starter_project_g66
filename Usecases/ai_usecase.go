package usecases

import (
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	"encoding/json"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type AIusecase struct {
// 	BlogUC BlogUseCase
// }

type AICommentUsecase struct {
	// AiComment domain.IAIInteraction
}

type AIFilterUsecase struct {
	BlogUC domain.IBlogUseCase
	// AiFilter domain.IAIInteraction
}

type AIBlogUsecase struct {
	// AiBlog domain.IAIInteraction
	BlogUC domain.IBlogUseCase
}

func NewAIusecaseComment() *AICommentUsecase {
	return &AICommentUsecase{
		// AiComment: newCommentinter,
	}
}
func NewAIusecaseBLog(bu domain.IBlogUseCase) *AIBlogUsecase {
	return &AIBlogUsecase{
		// AiBlog: newBlog,
		BlogUC: bu,
	}
}
func NewAIusecaseFilter(bu domain.IBlogUseCase) *AIFilterUsecase {
	return &AIFilterUsecase{
		BlogUC: bu,
		// AiFilter: newFliter,
	}
}

// func NewAIBlogUsecase(newbguc BlogUseCase,) *AIusecase {
// 	return &AIusecase{
// 		BlogUC:newbguc ,
// 	}
// }

func (au *AICommentUsecase) AICommentUsecase(userReq *domain.AICommentDTO, aIInteraction domain.IAIInteraction) (string, error) {
	res, err := aIInteraction.CallAIAndGetResponse("", userReq.UserMessage, userReq.Comment)
	log.Println("!!!!!!!!!!!!#in usecase", res)
	if err != nil {
		return "", err
	}
	if res.IsNilResponse {
		return "", errors.New("nil comment")
	}
	mainResponse := string(aIInteraction.ParseJsonBodyToDomain(res).(json.RawMessage))

	return mainResponse, nil
}

func (au *AIBlogUsecase) AIBlogUsecase(userID string, userReq *domain.AIBlogDTO, aIInteraction domain.IAIInteraction) (domain.Blog, error) {

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.Blog{}, err
	}
	blogDTO, err := au.BlogUC.GetBlogByID(objID)
	if err != nil {
		return domain.Blog{}, err
	}
	blogDomain := conv.ChangeToDomainBlog(blogDTO)
	jsonData, err := json.Marshal(blogDomain)
	if err != nil {
		return domain.Blog{}, err
	}

	res, err := aIInteraction.CallAIAndGetResponse("", "", string(jsonData))
	if err != nil {
		return domain.Blog{}, err
	}
	if res.IsNilResponse {
		return domain.Blog{}, errors.New("nil blog")
	}
	var blog *domain.Blog
	log.Println("💀 Before")
	blog = aIInteraction.ParseJsonBodyToDomain(res).(*domain.Blog)

	log.Println("💀 AFter")

	return *blog, nil
}

func (au *AIFilterUsecase) AIFilterUsecase(userReq *domain.AIBlogDTO, aIInteraction domain.IAIInteraction) ([]*domain.Blog, error) {
	res, err := aIInteraction.CallAIAndGetResponse("", userReq.UserMessage, "")
	if err != nil {
		return nil, err
	}
	resFilter := aIInteraction.ParseJsonBodyToDomain(res).(*domain.AIBlogFilter)
	log.Println("NICE", resFilter)
	resList, err := au.BlogUC.GetMainBlogByAIFitlter(resFilter)
	log.Println("Making the resList")
	if err != nil {
		return nil, err
	}
	resBlogList := []*domain.Blog{}
	log.Println("ምን ያደርጋል ኤአይ response is null")
	if resList == nil {
		return nil, nil
	}
	for _, each_blog_dto := range resList {
		resBlogList = append(resBlogList, conv.ChangeToDomainBlog(each_blog_dto))
	}
	return resBlogList, nil
}
