package domain

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthService interface {
	GenerateTokens(user *UserDTO) (string, string, error)
	// ValidateAccessToken(tokenStr string) (jwt.MapClaims, error)
	ValidateRefreshToken(tokenStr string) (string, error)
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
	OAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
}
type IAuthRepo interface {
	Save(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	Delete(token string) error
}

type IUserRepository interface { // eka was here
	Create(user *User) error

	// UpdatePassword(userEmail string, updatedPassword string) error
	// EditUserByEmail(userEmail string, updatedUserInfo *User) error
	FindByEmail(email string) (*UserDTO, error) //checks if user exisits or not
	UpdatePassword(userID, hashedPassword string) error
	CheckUserExistance(userEmail string) bool
	UpdateRole(email, role string) error
	UpdateUserByEmail(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	// DemoteUser(userEmail string) error
	// PromoteUser(userEmail string) error
	GetUserByID(userID string) (*UserDTO, error)
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
	Send(email string, token string) error
	GenerateRandomOTP() string
	SendResetLink(toEmail, subject, message string) error
}
type IBlogRepository interface {
	CreateBlog(blog *Blog, userID primitive.ObjectID) (*Blog,error)
	FindBlogByID(blogID primitive.ObjectID) (*BlogDTO, error)
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
	GetPopularityBlogByID(popValue primitive.ObjectID) (*PopularityDTO, error)
	CloseDataBase() error
}

type IBlogUseCase interface {
	CreateBlog(blog *Blog, userEmail string) (*Blog, error) //! Instead of userEmail as string we can pass userID instantly
	DeleteBlogByID(blogID string) error            // the controller will pass the a string from the url the usecase will change it to the objectID
	UpdateBlogByID(blogID string, updatedBlog *Blog) error
	GetBlogByID(blogID primitive.ObjectID) (*BlogDTO, error)
	// page number needed for the purpose of pagination
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*BlogDTO, error)
	LikeBlog(blogID primitive.ObjectID, userEmail string) error
	DisLikeBlog(blogID primitive.ObjectID, userEmail string) error
	CommentBlog(userEmail string, comment *CommentDTO, blogID primitive.ObjectID) error
	IncreaseView(blogID primitive.ObjectID) error
}

type IPasswordUsecase interface {
	GenerateResetToken(email string) error
	ResetPassword(token, newPassword string) error
}

type IAIInteraction interface {
	IsClientConnected() bool
	GenerateContent(prompt string) (*AIResponse, error)
	ParseJsonBodyToDomain(aiResponse *AIResponse) any
	CallAIAndGetResponse(developerMessage string, userMessage string, jsonBodyStirng string) (*AIResponse, error)
	IncrementInteractionCount()
	CloseAIConnection() error
}
