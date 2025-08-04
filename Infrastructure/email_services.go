package infrastructure

import (
	"fmt"
	"math/rand"
	"net/smtp"
)

type EmailService struct {
	From        string
	AppPass     string
	SmtpService string
	SmtPort     string
}

func NewOTP_service(from, app, smtpsev, smtport string) *EmailService {
	return &EmailService{
		From:        from,
		AppPass:     app,
		SmtpService: smtpsev,
		SmtPort:     smtport,
	}
}

func (s *EmailService) GenerateRandomOTP() string {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}
func (s *EmailService) Send(toEmail, otp string) error {

	auth := smtp.PlainAuth("", s.From, s.AppPass, s.SmtpService)
	msg := []byte(fmt.Sprintf("Subject: OTP Verification\n\nYour OTP is: %s", otp))
	return smtp.SendMail(s.SmtpService+":"+s.SmtPort, auth, s.From, []string{toEmail}, msg)
}
