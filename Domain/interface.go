package domain


type IUserRepository interface {

import (
	"blog_starter_project_g66/Delivery/controllers"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserRepository interface { // eka was here
	Create(user *User) error
	FetchByEmail(userEmail string) (*User, error) //checks if user exisits or not
	CheckUserExistance(userEmail string) bool
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
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*controllers.BlogDTO, error)
	CheckBlogExistance(blogID primitive.ObjectID) bool
	CloseDataBase() error

}

type IBlogUseCase interface {
	CreateBlog(blog *Blog, userEmail string) error
	DeleteBlogByID(blogID string) error // the controller will pass the a string from the url the usecase will change it to the objectID
	UpdateBlogByID(blogID string, updatedBlog *Blog) error
	// page number needed for the purpose of pagination
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*controllers.BlogDTO, error)
}
