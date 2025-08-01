package infrastructure

import (
	"blog_starter_project_g66/Domain"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type OTP_service struct{
	otpService *domain.GenerateOTP
}

func NewOTP_service() *OTP_service{
	return &OTP_service{}
}

func GenerateRandomOTP() string {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}
func (s *OTP_service) Send(toEmail, otp string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	from := os.Getenv("FROM")
	appPass := os.Getenv("APPPASS")
	smtpServer := os.Getenv("SMTPSERVER")
	smtpPort := os.Getenv("SMTPPORT")
	auth := smtp.PlainAuth("", from, appPass, smtpServer)
	msg := []byte(fmt.Sprintf("Subject: OTP Verification\n\nYour OTP is: %s", otp))
	return smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{toEmail}, msg)
}