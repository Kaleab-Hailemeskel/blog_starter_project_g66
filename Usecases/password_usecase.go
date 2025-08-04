package usecases

import (
	"blog_starter_project_g66/Domain"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type passwordUsecase struct {
	userRepo 		domain.IUserRepository
	emailService 	domain.IEmailService
	jwtSecret 		string 
}

func NewPasswordUsecase(repo domain.IUserRepository, emailServ domain.IEmailService, jwtSec string) domain.IPasswordUsecase {
	return &passwordUsecase{
		userRepo: repo,
		emailService: emailServ,
		jwtSecret: jwtSec,
	}
}

func (u *passwordUsecase) GenerateResetToken(email string) error {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	claims := jwt.MapClaims{
		"user_id": 	user.UserID.Hex(),
		"email":	user.Email,
		"exp":		time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.jwtSecret))

	if err != nil {
		return err
	}

	return u.emailService.SendPasswordReset(user.Email, signedToken)
}

func (u *passwordUsecase) ResetPassword(tokenStr, newPassword string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("invalid email in token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return u.userRepo.UpdatePassword(email, string(hashedPassword))

}