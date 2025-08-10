package main

import (
	"blog_starter_project_g66/Delivery/controllers"
	"blog_starter_project_g66/Delivery/oauth"
	"blog_starter_project_g66/Delivery/routers"
	infrastructure "blog_starter_project_g66/Infrastructure"
	repositories "blog_starter_project_g66/Repositories"
	usecases "blog_starter_project_g66/Usecases"
	"blog_starter_project_g66/config"
	"log"
)

func main() {
	config.InitEnv()
	mongoClient, err := repositories.Connect()
	if err != nil {
		log.Fatal("❌Failed to connect:", err)
	}
	defer mongoClient.Disconnect()

	from := config.FROM
	appPass := config.APPPASS
	smtpServer := config.SMTPSERVER
	smtpPort := config.SMTPPORT
	user := config.SMTPUSER
	jwtSecret := config.JWTSECRET

	oauth.InitOAuth()
	// Initialize DataBase Repository
	userRepo := repositories.NewUserRepository()
	blogRepo := repositories.NewBlogDataBaseService()
	popularityRepo := repositories.NewBlogPopularityDataBaseService()
	authRepo := repositories.NewRefreshTokenRepository()
	otpService := repositories.NewUserOTPRepository()

	// Initialize Services
	passwaordService := infrastructure.NewPasswordService()
	authMiddleware := infrastructure.NewAuthMiddleware(authRepo)
	authService := infrastructure.NewJWTService(authRepo)
	emailService := infrastructure.NewOTP_service(from, appPass, smtpServer, smtpPort, user)

	// Initialize UseCases
	oauthUsecase := usecases.NewOAuthUsecase(userRepo, authService)
	userUsecase := usecases.NewUserUsecase(userRepo, passwaordService, otpService, emailService, authService, authRepo)
	blogUsecase := usecases.NewBlogUseCase(blogRepo, userRepo, popularityRepo)
	passwordUsecase := usecases.NewPasswordUsecase(userRepo, emailService, jwtSecret)

	// Initialize DataBase Repository
	userController := controllers.NewUserUsecase(userUsecase, oauthUsecase)
	passwordController := controllers.NewPasswordController(passwordUsecase)
	blogController := controllers.NewController(blogUsecase, userUsecase)

	//Initialize AI Interaction
	aiCommentUsecase := usecases.NewAIusecaseComment()
	aiBlogUsecase := usecases.NewAIusecaseBLog(blogUsecase)
	aiFliterUsecase := usecases.NewAIusecaseFilter()
	aiController := controllers.NewAIController(aiCommentUsecase, aiBlogUsecase, aiFliterUsecase)

	// Initialize Routers
	routers.Router(userController, passwordController, blogController, authMiddleware, aiController)

}
