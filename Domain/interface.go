package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserRepository interface { // eka was here
	Create(user *User) error
	EditUserByEmail(userEmail string, updatedUserInfo *User) error
	CheckUserExistance(userEmail string) bool
	FindByEmail(email string) (*UserDTO, error) //checks if user exisits or not
	UpdatePassword(userID, hashedPassword string) error
	CloseDataBase() error
}

type IUserValidation interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) bool
	Hashpassword(password string) string
	ComparePassword(userPassword, password string) error
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

type IBlogUseCase interface {
	CreateBlog(blog *Blog, userEmail string) error
	DeleteBlogByID(blogID string) error // the controller will pass the a string from the url the usecase will change it to the objectID
	UpdateBlogByID(blogID string, updatedBlog *Blog) error
	// page number needed for the purpose of pagination
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*BlogDTO, error)
}

type IEmailService interface {
	SendPasswordReset(email string, token string) error
}

type IPasswordUsecase interface {
	GenerateResetToken(email string) error
	ResetPassword(token, newPassword string) error
}