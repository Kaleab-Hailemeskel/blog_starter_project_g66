package infrastructure

import (
	"blog_starter_project_g66/Domain"
	"fmt"

	"gopkg.in/mail.v2"
)

type SMTPService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPEmailservice(host string, port int, username, password, from string) domain.IEmailService{
	return &SMTPService{
		Host:     host,
        Port:     port,
        Username: username,
        Password: password, 
        From:     from,
	}
} 

func (s *SMTPService) SendPasswordReset(email , token string) error{
	link := "http://localhost:3000/reset-password?token=" + token
	m := mail.NewMessage()
	m.SetHeader("From", s.From) 
    m.SetHeader("To", email)
    m.SetHeader("Subject", "Password Reset Request")
    m.SetBody("text/plain", fmt.Sprintf("Click here to reset your password: %s", link))

    d := mail.NewDialer(s.Host, s.Port, s.Username, s.Password)
    return d.DialAndSend(m)
} 