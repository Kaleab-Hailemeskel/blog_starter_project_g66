package domain

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DATABASE Repositorys
type IUserRepository interface {
	// eka was here
	Create(user *User) error
	FindByEmail(email string) (*UserDTO, error) //checks if user exisits or not
	UpdatePassword(userID, hashedPassword string) error
	CheckUserExistance(userEmail string) bool
	UpdateRole(email, role string) error
	UpdateUserByEmail(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	GetUserByID(userID string) (*UserDTO, error)
	CreateSuperAdmin() error
	CloseDataBase() error
}
type IAuthRepo interface {
	Save(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	Delete(token string) error
	CloseDataBase() error
}
type IUserOTPRepository interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
	CloseDataBase() error
}
type IBlogRepository interface {
	IsClientConnected() bool // just for testing on the testify purpose
	CreateBlog(blog *Blog, userID primitive.ObjectID) (*BlogDTO, error)
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
	DeletePopularityBlogByID(blogID primitive.ObjectID) error
	UserDisLikeBlogByID(blogID primitive.ObjectID, userID primitive.ObjectID, revert bool) error
	CreateBlogPopularity(blogID primitive.ObjectID) (*PopularityDTO, error)
	UpdatePopularityValueByBlogID(blogID primitive.ObjectID, calculatedValue int) error
	CommentBlogByID(blogID primitive.ObjectID, commentDTO *CommentDTO) error
	IncreaseBlogViewByID(blogID primitive.ObjectID) error
	BlogPostViewCountByID(blogID primitive.ObjectID) (int, error)
	BlogPostPopularityValueByID(blogID primitive.ObjectID) (int, error)
	BlogPostLikeCountByID(blogID primitive.ObjectID) (int, error)
	BlogPostDisLikeCountByID(blogID primitive.ObjectID) (int, error)
	BlogPostCommentCountByID(blogID primitive.ObjectID) (int, error)
	GetPopularityBlogByID(blogID primitive.ObjectID) (*PopularityDTO, error)
	CloseDataBase() error
}

// SERVICE implementation
type IAuthService interface {
	GenerateTokens(user *UserDTO) (string, string, error)
	ValidateRefreshToken(tokenStr string) (string, error)
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
	OAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
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

type IAIInteraction interface {
	IsClientConnected() bool
	GenerateContent(prompt string) (*AIResponse, error)
	ParseJsonBodyToDomain(aiResponse *AIResponse) any
	CallAIAndGetResponse(developerMessage string, userMessage string, jsonBodyStirng string) (*AIResponse, error)
	IncrementInteractionCount()
	CloseAIConnection() error
}

// USECASE declarations
type IUserUseCase interface {
	HandleRegistration(user *User) error
	SendOTP(user *User) error
	VerifyOTP(email, otp string) (bool, error)
	PromoteUser(actor, target string) error
	DemoteUser(actor, target string) error
	Login(email, password string) (*AuthTokens, error)
	Refresh(oldRefreshToken string) (*AuthTokens, error)
	Logout(refreshToken string) error
	UpdateProfile(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	GetUserByEmail(email string) (*UserDTO, error)
}
type IBlogUseCase interface {
	CreateBlog(blog *Blog, userEmail string) (*BlogDTO, error)
	DeleteBlogByID(blogID string) error
	UpdateBlogByID(blogID string, updatedBlog *Blog) error
	GetBlogByID(blogID primitive.ObjectID) (*BlogDTO, error)
	// page number needed for the purpose of pagination
	GetAllBlogsByFilter(url_filter *Filter, pageNumber int) ([]*BlogDTO, error)
	LikeBlog(blogID primitive.ObjectID, userEmail string) error
	DisLikeBlog(blogID primitive.ObjectID, userEmail string) error
	CommentBlog(userEmail string, comment *CommentDTO, blogID primitive.ObjectID) error
	IncreaseView(blogID primitive.ObjectID) error

	GetPopularityBlogByID(blogID primitive.ObjectID) (*PopularityDTO, error)
	CalcualtePopularity(blog *PopularityDTO) int
	CommentBlogByID(blogID primitive.ObjectID, commentDTO *Comment) error

	GetMainBlogAndPopularityBlogByID(blogID primitive.ObjectID) (*BlogDTO, *PopularityDTO, error)

	GetMainBlogByAIFitlter(aiFilter *AIBlogFilter) ([]*BlogDTO, error)
}
type IPasswordUsecase interface {
	GenerateResetToken(email string) error
	ResetPassword(token, newPassword string) error
}
type IOAuthUsecase interface {
	HandleOAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
}

type IAICommentUsecase interface {
	AICommentUsecase(userReq *AICommentDTO, aIInteraction IAIInteraction) (string, error)
}
type IAIBlogUsecase interface {
	AIBlogUsecase(userID string, userReq *AIBlogDTO, aIInteraction IAIInteraction) (Blog, error)
}
type IAIFilterUsecase interface {
	AIFilterUsecase(*AIBlogDTO, IAIInteraction) ([]*Blog, error)
}
