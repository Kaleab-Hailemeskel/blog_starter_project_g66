package usecases

import (
	domain "blog_starter_project_g66/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUseCase struct {
	BlogDataBase       domain.IBlogRepository
	UserDataBase       domain.IUserRepository
	PopularityDataBase domain.IPopularityRepository
}

func NewBlogUseCase(blogRepo domain.IBlogRepository, userRepo domain.IUserRepository, PopRepo domain.IPopularityRepository) *BlogUseCase {
	return &BlogUseCase{
		BlogDataBase:       blogRepo,
		UserDataBase:       userRepo,
		PopularityDataBase: PopRepo,
	}
}
func (bluc *BlogUseCase) CreateBlog(blog *domain.Blog, ownerEmail string) error {
	user, err := bluc.UserDataBase.FindByEmail(ownerEmail)
	if err != nil {
		return err
	}
	userID := user.UserID
	return bluc.BlogDataBase.CreateBlog(blog, userID)
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
	if pageNumber < 1 {
		pageNumber = 1
	}
	return bluc.BlogDataBase.GetAllBlogsByFilter(url_filter, pageNumber)
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
	return blue.PopularityDataBase.UserLikeBlogByID(blogID, userDTO.UserID, false)
}
func (blue *BlogUseCase) DisLikeBlog(blogID primitive.ObjectID, userEmail string) error {
	userDTO, err := blue.UserDataBase.FindByEmail(userEmail)
	if err != nil {
		return err
	}
	if blue.PopularityDataBase.CheckUserLikeBlogID(blogID, userDTO.UserID) {
		err = blue.PopularityDataBase.UserLikeBlogByID(blogID, userDTO.UserID, true)
		if err != nil {
			return err
		}
	}
	return blue.PopularityDataBase.UserDisLikeBlogByID(blogID, userDTO.UserID, false)
}
func (blue *BlogUseCase) CommentBlog(userEmail string, comment *domain.CommentDTO, blogID primitive.ObjectID) error {
	userDTO, err := blue.UserDataBase.FindByEmail(userEmail)
	if err != nil {
		return err
	}
	// setting name and ID on the comment so that the commentBlogByID will track it and update the comment
	comment.UserID = userDTO.UserID
	comment.UserName = userDTO.UserName
	return blue.PopularityDataBase.CommentBlogByID(blogID, comment)
}
func (blue *BlogUseCase) IncreaseView(blogID primitive.ObjectID) error { // we don't have to care about who watch the blog, just add the view_count.
	return blue.PopularityDataBase.IncreaseBlogViewByID(blogID)
}
