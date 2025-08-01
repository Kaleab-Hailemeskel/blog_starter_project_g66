package usecases

import (
	domain "blog_starter_project_g66/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUseCase struct {
	BlogDataBase domain.IBlogRepository
	UserDataBase domain.IUserRepository
}

func NewBlogUseCase(blogRepo domain.IBlogRepository, userRepo domain.IUserRepository) domain.IBlogUseCase {
	return &BlogUseCase{
		BlogDataBase: blogRepo,
		UserDataBase: userRepo,
	}
}
func (bluc *BlogUseCase) CreateBlog(blog *domain.Blog, ownerEmail string) error {
	user, err := bluc.UserDataBase.FetchByEmail(ownerEmail)
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
