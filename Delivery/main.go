package main

import (
	"blog_starter_project_g66/Delivery/controllers"
	"blog_starter_project_g66/Delivery/routers"
	"blog_starter_project_g66/Infrastructure"
	"blog_starter_project_g66/Repositories"
	"blog_starter_project_g66/Usecases"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	mongoURI := os.Getenv("MONGO_CONNECTION_STRING")
	jwtSecret := os.Getenv("JWT_SECERT")

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")


	mongoClient, err := repositories.Connect()
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	defer mongoClient.Disconnect()

	userRepo := repositories.INewUserRepository(mongoClient)
	emailService := infrastructure.NewSMTPEmailservice(host, port, user, pass, from)
	passwordUsecase := usecases.NewPasswordUsecase(userRepo, emailService, jwtSecret)
	passwordController := controllers.NewPasswordController(passwordUsecase)
	// routers.Router(,pc)



}
