package domain

type UserRepository interface {
	Create(user *User) error
	Login(email, password string) (*User, error) //authenticate and return user
	FetchByEmail(email string) (bool, error)     //checks if user exisits or not
}

type IUserRepository interface {
	CreateBlog(blog *Blog) error
	DeleteBlogByID(blogID string) error
	UpdateBlogByID(blogID string, updatedBlog *Blog) error
	GetAllBlogs(pageNumber int) (*[]Blog, error) // i used pageNumber for pagination so that large contents accessed by page number for the firstRequest request, after that the user will have the hashed indexed, and will be able to access next databases after the lastIndexBlog
	GetAllBlogsByFilter(filter Filter, pageNumber int) (*[]Blog, error)
	CheckBlogExistance(blogID string) bool
	CloseDataBase() error
}
