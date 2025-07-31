package domain

import "time"

type UserRepository interface {
	Create(user *User) error
	Login(email, password string) (*User, error) //authenticate and return user
	FetchByEmail(email string) (bool, error)     //checks if user exisits or not
}

type IUserRepository interface {
	CreateBlog(blog *Blog) error
	DeleteBlogByID(blogID string) error
	UpdateBlog(blogID string, updatedBlog *Blog) error
	GetAllBlogs(indexBlogID *string) (*[]Blog, error) // i used indexBlog ID for pagination puposes, and *string arg is needed to pass nil for the firstRequest request, after that the user will have the hashed indexed, and will be able to access next databases after the lastIndexBlog
	GetAllBlogsByFilter(tagFilter string, dateFilter time.Time, popularityValue int, indexBlog *string) (*[]Blog, error)
	GetByID(blogID string) (*Blog, error)
	CheckBlogExistance(blogID string) bool
	CloseDataBase() error
}