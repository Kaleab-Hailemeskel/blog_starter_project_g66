package main

import (
	"blog_starter_project_g66/Delivery/controllers"
	"blog_starter_project_g66/Delivery/routers"
	"blog_starter_project_g66/Infrastructure"
	"blog_starter_project_g66/Repositories"
	"blog_starter_project_g66/Usecases"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	mongoClient, err := repositories.Connect()
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer mongoClient.Disconnect()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	from := os.Getenv("FROM")
	appPass := os.Getenv("APPPASS")
	smtpServer := os.Getenv("SMTPSERVER")
	smtpPort := os.Getenv("SMTPPORT")

	authRepo := repositories.NewRefreshTokenRepository(mongoClient)
	authService := infrastructure.NewJWTService()
	emailService := infrastructure.NewOTP_service(from, appPass, smtpServer, smtpPort)
	userRepo := repositories.NewUserRepository(mongoClient)
	otpService := repositories.NewUserOTPRepository(mongoClient)
	passwaordService := infrastructure.NewPasswordService()
	userUsecase := usecases.NewUserUsecase(userRepo, passwaordService, otpService, emailService,authService,authRepo)
	userController := controllers.NewUserUsecase(userUsecase)

	routers.Router(userController)
}
