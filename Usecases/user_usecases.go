package usecases

import (
	"blog_starter_project_g66/Domain"
	"errors"
	"time"
)

type UserUsecase struct {
	userinterface domain.IUserRepository
	userVaildate domain.IUserValidation
	userOTP domain.IUserOTP 
	generateotp domain.IEmailService

}

func NewUserUsecase(ui domain.IUserRepository,uv domain.IUserValidation, uo domain.IUserOTP, emailService domain.IEmailService) *UserUsecase{
	return &UserUsecase{
		userinterface: ui,
		userVaildate: uv,
		userOTP: uo,
		generateotp: emailService,
	}
}

func (uc *UserUsecase) HandleRegistration(user *domain.User) error {
	existing := uc.userinterface.CheckUserExistance(user.Email)

	if existing {
		return errors.New("user already exists")
	} 

	isvaild_email := uc.userVaildate.IsValidEmail(user.Email)
	ispassword_strong := uc.userVaildate.IsStrongPassword(user.Password)

	if !ispassword_strong || !isvaild_email {
		return errors.New("invalid password or email")
	}

	hashpassword := uc.userVaildate.Hashpassword(user.Password)
	user.Password = hashpassword

	

	err := uc.SendOTP(user)

	if err != nil {
		return errors.New("failed to send OTP: " + err.Error())
	}

	return nil
}
func (uc *UserUsecase) SendOTP(user *domain.User) error {
	otp := uc.generateotp.GenerateRandomOTP()
	entry := domain.UserUnverified{
		UserName: user.UserName ,
		Email:     user.Email,
		OTP:       otp,
		Password: user.Password,
		Role: user.Role,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := uc.generateotp.Send(user.Email, otp); err != nil {
		return err
	}
	return uc.userOTP.StoreOTP(entry)
}

func (uc *UserUsecase) VerifyOTP(email, otp string) (bool, error) {
	entry, err := uc.userOTP.FindOTP(email)
	if err != nil || entry == nil {
		return false, err
	}
	if time.Now().After(entry.ExpiresAt) || entry.OTP != otp {
		return false, nil
	}
	verifiedUser := &domain.User{
	UserName:       entry.UserName,
	Email:          entry.Email,
	Password:       entry.Password,
	Role:           entry.Role,
}
err = uc.userinterface.Create(verifiedUser)
	
	if err != nil {
		return false,err
	}
	_ = uc.userOTP.DeleteOTP(email)
	return true, nil
}
