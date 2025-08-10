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

	"github.com/joho/godotenv"
)

func main() {
	config.InitEnv()
	mongoClient, err := repositories.Connect()
	if err != nil {
		log.Fatal("❌Failed to connect:", err)
	}
	defer mongoClient.Disconnect()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	from := config.FROM
	appPass := config.APPPASS
	smtpServer := config.SMTPSERVER
	smtpPort := config.SMTPPORT
	user := config.SMTPUSER
	jwtSecret := config.JWTSECRET

	oauth.InitOAuth()

	authRepo := repositories.NewRefreshTokenRepository(mongoClient)
	authMiddleware := infrastructure.NewAuthMiddleware(authRepo)
	authService := infrastructure.NewJWTService(authRepo)
	emailService := infrastructure.NewOTP_service(from, appPass, smtpServer, smtpPort, user)
	userRepo := repositories.NewUserRepository()
	oauthUsecase := usecases.NewOAuthUsecase(userRepo, authService)
	otpService := repositories.NewUserOTPRepository(mongoClient)
	passwaordService := infrastructure.NewPasswordService()
	userUsecase := usecases.NewUserUsecase(userRepo, passwaordService, otpService, emailService, authService, authRepo)
	userController := controllers.NewUserUsecase(userUsecase, oauthUsecase)

	passwordUsecase := usecases.NewPasswordUsecase(userRepo, emailService, jwtSecret)
	passwordController := controllers.NewPasswordController(passwordUsecase)

	blogRepo := repositories.NewBlogDataBaseService()
	popularityRepo := repositories.NewBlogPopularityDataBaseService()
	blogUsecase := usecases.NewBlogUseCase(blogRepo, userRepo, popularityRepo)
	blogController := controllers.NewController(blogUsecase)

	aiCommentUsecase :=usecases.NewAIusecaseComment()
	aiBlogUsecase := usecases.NewAIusecaseBLog(blogUsecase)
	aiFliterUsecase := usecases.NewAIusecaseFilter()
	aiController := controllers.NewAIController(aiCommentUsecase,aiBlogUsecase,aiFliterUsecase)
	routers.Router(userController, passwordController, blogController, authMiddleware, aiController)

}
