package usecases

import (
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// Recommended Value 0.6 to 0.8
	WEIGHT_1 = 0.7
	// Recommended Value 0.2 to 0.4
	WEIGHT_2 = 0.3
)

type BlogUseCase struct {
	BlogDataBase       domain.IBlogRepository
	UserDataBase       domain.IUserRepository
	PopularityDataBase domain.IPopularityRepository
}

func NewBlogUseCase(blogRepo domain.IBlogRepository, userRepo domain.IUserRepository, PopRepo domain.IPopularityRepository) domain.IBlogUseCase {
	return &BlogUseCase{
		BlogDataBase:       blogRepo,
		UserDataBase:       userRepo,
		PopularityDataBase: PopRepo,
	}
}
func (bluc *BlogUseCase) CalcualtePopularity(blog *domain.PopularityDTO) int {

	numOfComment := len(blog.Comments)
	numOfLike := len(blog.Likes)
	numOfDisLike := len(blog.Likes)
	numOfView := blog.ViewCount

	CalcPopularity := float32(1)
	SentimentScore := float32(0)
	if numOfLike+numOfDisLike != 0 {
		SentimentScore = float32(numOfLike / (numOfLike + numOfDisLike))
	}
	EngagementScore := float32(0)
	if numOfView != 0 {
		EngagementScore = float32((numOfLike + numOfDisLike + numOfComment) / numOfView)
	}
	CalcPopularity = 100 * ((WEIGHT_1 * SentimentScore) + (WEIGHT_2 * EngagementScore))

	return int(CalcPopularity)
}

func (bluc *BlogUseCase) CommentBlogByID(blogID primitive.ObjectID, comment *domain.Comment) error {
	return bluc.PopularityDataBase.CommentBlogByID(blogID, conv.ChangeToDTOComment(comment))
}
func (bluc *BlogUseCase) CreateBlog(blog *domain.Blog, ownerEmail string) (*domain.BlogDTO, error) {
	user, err := bluc.UserDataBase.FindByEmail(ownerEmail)
	if err != nil {
		return nil, err
	}
	userID := user.UserID
	createdBlog, err := bluc.BlogDataBase.CreateBlog(blog, userID)
	_, poplarityError := bluc.PopularityDataBase.CreateBlogPopularity(createdBlog.BlogID)
	if err == nil && poplarityError != nil {
		return nil, fmt.Errorf("_Main Blog Created But Popularity is not >> 💀💀💀❌❌❌")
	}
	return createdBlog, err
}
func (bluc *BlogUseCase) DeleteBlogByID(blogID string) error {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	return bluc.BlogDataBase.DeleteBlogByID(blogObjID)
}
func (bluc *BlogUseCase) UpdateBlogByID(blogID string, updatedBlog *domain.Blog) error {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	return bluc.BlogDataBase.UpdateBlogByID(blogObjID, updatedBlog)
}
func (bluc *BlogUseCase) GetAllBlogsByFilter(url_filter *domain.Filter, pageNumber int) ([]*domain.BlogDTO, error) {
	if url_filter == nil {
		return nil, fmt.Errorf("filter cannot be nil")
	}
	if pageNumber < 1 {
		pageNumber = 1
	}
	if res, err := bluc.BlogDataBase.GetAllBlogsByFilter(url_filter, pageNumber); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
func (bluc *BlogUseCase) GetBlogByID(blogID primitive.ObjectID) (*domain.BlogDTO, error) {
	return bluc.BlogDataBase.FindBlogByID(blogID)
}
func (blue *BlogUseCase) LikeBlog(blogID primitive.ObjectID, userEmail string) error {
	userDTO, err := blue.UserDataBase.FindByEmail(userEmail)
	if err != nil {
		return err
	}
	if blue.PopularityDataBase.CheckUserDisLikeBlogID(blogID, userDTO.UserID) {
		err = blue.PopularityDataBase.UserDisLikeBlogByID(blogID, userDTO.UserID, true)
		if err != nil {
			return err
		}
	}
	revertIt := blue.PopularityDataBase.CheckUserLikeBlogID(blogID, userDTO.UserID)
	retErr := blue.PopularityDataBase.UserLikeBlogByID(blogID, userDTO.UserID, revertIt)
	recalcError := blue.RecalcuatePopularityValue(&blogID, nil)
	if retErr == nil && recalcError != nil {
		return recalcError
	}
	return retErr
}
func (blue *BlogUseCase) DisLikeBlog(blogID primitive.ObjectID, userEmail string) error { //! go routine zone
	userDTO, err := blue.UserDataBase.FindByEmail(userEmail)
	if err != nil {
		return err
	}
	if blue.PopularityDataBase.CheckUserLikeBlogID(blogID, userDTO.UserID) {
		err = blue.PopularityDataBase.UserLikeBlogByID(blogID, userDTO.UserID, true) // this func is independent on it's own so use go rountine. LIKEuSEcASEsHOULDaLSOhAVEiT. <- real msg by the way 😅
		if err != nil {
			return err
		}
	}
	revertIt := blue.PopularityDataBase.CheckUserDisLikeBlogID(blogID, userDTO.UserID)
	retErr := blue.PopularityDataBase.UserDisLikeBlogByID(blogID, userDTO.UserID, revertIt)
	recalcError := blue.RecalcuatePopularityValue(&blogID, nil)
	if retErr == nil && recalcError != nil {
		return recalcError
	}
	return retErr
}
func (blue *BlogUseCase) CommentBlog(userEmail string, comment *domain.CommentDTO, blogID primitive.ObjectID) error {
	userDTO, err := blue.UserDataBase.FindByEmail(userEmail)
	if err != nil {
		return err
	}
	// setting name and ID on the comment so that the commentBlogByID will track it and update the comment
	comment.OwnerID = userDTO.UserID
	comment.UserName = userDTO.UserName
	commentErr := blue.PopularityDataBase.CommentBlogByID(blogID, comment)
	if commentErr == nil {
		if calculateErr := blue.RecalcuatePopularityValue(&blogID, nil); calculateErr != nil {
			return calculateErr
		}
	}
	return commentErr
}
func (blue *BlogUseCase) IncreaseView(blogID primitive.ObjectID) error { // we don't have to care about who watch the blog, just add the view_count.
	increaseViewError := blue.PopularityDataBase.IncreaseBlogViewByID(blogID)
	if increaseViewError != nil {
		return increaseViewError
	}
	anotherErr := blue.RecalcuatePopularityValue(&blogID, nil)
	if anotherErr != nil {
		return anotherErr
	}
	return increaseViewError
}
func (blue *BlogUseCase) GetPopularityBlogByID(blogID primitive.ObjectID) (*domain.PopularityDTO, error) {
	return blue.PopularityDataBase.GetPopularityBlogByID(blogID)
}

func (blue *BlogUseCase) GetMainBlogAndPopularityBlogByID(blogID primitive.ObjectID) (*domain.BlogDTO, *domain.PopularityDTO, error) {
	//! go routine to not wait until the first database query finishes
	blue.IncreaseView(blogID) // go routine maybe
	resBlog, errBlog := blue.GetBlogByID(blogID)
	if errBlog != nil {
		return nil, nil, errBlog
	}
	resBlogPopularity, errPopularity := blue.GetPopularityBlogByID(blogID)
	if errPopularity != nil {
		return nil, nil, errPopularity
	}
	return resBlog, resBlogPopularity, nil
}

// Either pass the blogID or the popularityBlog, but not both
func (blue *BlogUseCase) RecalcuatePopularityValue(blogID *primitive.ObjectID, popularityBlog *domain.PopularityDTO) error {
	// ei
	if blogID != nil { // if blogID present
		popBlogRes, err := blue.GetPopularityBlogByID(*blogID)
		if err != nil {
			return err
		}
		popVal := blue.CalcualtePopularity(popBlogRes)
		return blue.PopularityDataBase.UpdatePopularityValueByBlogID(*blogID, popVal)
	} else if popularityBlog != nil { // if popularityBlog is passed
		popVal := blue.CalcualtePopularity(popularityBlog)
		return blue.PopularityDataBase.UpdatePopularityValueByBlogID(popularityBlog.BlogID, popVal)
	}
	return fmt.Errorf("recalculating the PopularityValue in null values doesn't make sense")
}
