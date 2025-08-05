package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type IUserRepository interface { // eka was here
	Create(user *User) error
	// FetchByEmail(userEmail string) (*UserDTO, error) //checks if user exisits or not
	// UpdatePassword(userEmail string, updatedPassword string) error
	// EditUserByEmail(userEmail string, updatedUserInfo *User) error
	FindByEmail(email string) (*UserDTO, error) //checks if user exisits or not
	UpdatePassword(userID, hashedPassword string) error
	CheckUserExistance(userEmail string) bool
	DemoteUser(userEmail string) error
	PromoteUser(userEmail string) error
	CloseDataBase() error
}

type IUserValidation interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) bool
	Hashpassword(password string) string
	ComparePassword(userPassword, password string) error
}
type IUserOTP interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
}
type IEmailService interface {
    Send( email string, token string) error
	// SendPasswordReset(to string, subject string, body string) error
	GenerateRandomOTP() string 
	
}
type IBlogRepository interface {
	CreateBlog(blog *Blog, userID primitive.ObjectID) error
	FindBlogByID(blogID primitive.ObjectID) (*Blog, error)
	DeleteBlogByID(blogID primitive.ObjectID) error
	UpdateBlogByID(blogID primitive.ObjectID, updatedBlog *Blog) error
	// page number needed for the purpose of pagination
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*BlogDTO, error)
	CheckBlogExistance(blogID primitive.ObjectID) bool
	CloseDataBase() error
}
type IPopularityRepository interface {
	CheckUserLikeBlogID(blogID primitive.ObjectID, userID primitive.ObjectID) bool
	CheckUserDisLikeBlogID(blogID primitive.ObjectID, userID primitive.ObjectID) bool
	UserLikeBlogByID(blogID primitive.ObjectID, userID primitive.ObjectID, revert bool) error // revert boolean helps to undo the like while disliking the blog
	UserDisLikeBlogByID(blogID primitive.ObjectID, userID primitive.ObjectID, revert bool) error
	CreateBlogPopularity(blogID primitive.ObjectID) error
	CommentBlogByID(blogID primitive.ObjectID, commentDTO *CommentDTO) error
	IncreaseBlogViewByID(blogID primitive.ObjectID) error
	BlogPostLikeCountByID(blogID primitive.ObjectID) (int, error)
	BlogPostDisLikeCountByID(blogID primitive.ObjectID) (int, error)
	BlogPostCommentCountByID(blogID primitive.ObjectID) (int, error)
	CloseDataBase() error
}

type IBlogUseCase interface {
	CreateBlog(blog *Blog, userEmail string) error //! Instead of userEmail as string we can pass userID instantly
	DeleteBlogByID(blogID string) error            // the controller will pass the a string from the url the usecase will change it to the objectID
	UpdateBlogByID(blogID string, updatedBlog *Blog) error
	// page number needed for the purpose of pagination
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*BlogDTO, error)
}

type IPasswordUsecase interface {
	GenerateResetToken(email string) (string, error)
	ResetPassword(token, newPassword string) error
}
