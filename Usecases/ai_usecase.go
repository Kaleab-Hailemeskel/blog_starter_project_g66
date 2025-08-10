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

type AICommentUsecase struct{
	// AiComment domain.IAIInteraction
}

type AIFilterUsecase struct{
	
	// AiFilter domain.IAIInteraction
}

type AIBlogUsecase struct{
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
func NewAIusecaseFilter() *AIFilterUsecase{
	return &AIFilterUsecase{
		// AiFilter: newFliter,
	}
}
// func NewAIBlogUsecase(newbguc BlogUseCase,) *AIusecase {
// 	return &AIusecase{
// 		BlogUC:newbguc ,
// 	}
// }

func (au *AICommentUsecase) AICommentUsecase(userReq *domain.AICommentDTO,aIInteraction domain.IAIInteraction) (string, error) {
	res , err :=aIInteraction.CallAIAndGetResponse("", userReq.UserMessage,userReq.Comment)
	log.Println("!!!!!!!!!!!!#in usecase",res)
	if err != nil{
		return "",err
	}
	if res.IsNilResponse{
		return "",errors.New("nil comment")
	}
	mainResponse := string(aIInteraction.ParseJsonBodyToDomain(res).(json.RawMessage))

	return mainResponse, nil
}

func (au *AIBlogUsecase) AIBlogUsecase(userID string,userReq *domain.AIBlogDTO, aIInteraction domain.IAIInteraction) (domain.Blog, error){
	

	objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return domain.Blog{},err
    }
	blogDTO, err:=au.BlogUC.GetBlogByID(objID)
	if err != nil{
		return domain.Blog{},err
	}
	blogDomain :=conv.ChangeToDomainBlog(blogDTO)
	jsonData, err := json.Marshal(blogDomain)
	if err != nil{
		return domain.Blog{},err
	}

	res , err :=aIInteraction.CallAIAndGetResponse("","",string(jsonData))
	if err != nil{
		return domain.Blog{},err
	}
	if res.IsNilResponse{
		return domain.Blog{},errors.New("nil blog")
	}
	var blog domain.Blog
	err = json.Unmarshal(aIInteraction.ParseJsonBodyToDomain(res).(json.RawMessage), &blog)
	if err != nil {

		return domain.Blog{},err
	}

	return blog, nil
}

func (au *AIFilterUsecase) AIFilterUsecase(userReq *domain.AIBlogDTO,aIInteraction domain.IAIInteraction) (domain.Blog, error){
	var blog domain.Blog
	res, err := aIInteraction.CallAIAndGetResponse("",userReq.UserMessage,"")
	if err != nil {

		return domain.Blog{},err
	}
	err = json.Unmarshal(aIInteraction.ParseJsonBodyToDomain(res).(json.RawMessage), &blog)
	if err != nil {

		return domain.Blog{},err
	}


	return blog, nil
}