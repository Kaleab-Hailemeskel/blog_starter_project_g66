package usecases

import (
	domain "blog_starter_project_g66/Domain"
	"errors"
	"time"
)

type UserUsecase struct {
	userinterface domain.IUserRepository
	userVaildate domain.IUserValidation
	userOTP domain.IUserOTP 
	generateotp domain.GenerateOTP

}

func NewUserUsecase(ui domain.IUserRepository,uv domain.IUserValidation, uo domain.IUserOTP) *UserUsecase{
	return &UserUsecase{
		userinterface: ui,
		userVaildate: uv,
		userOTP: uo,
	}
}

func (uc *UserUsecase) HandleRegistration(user *domain.User) error {
	existing := uc.userinterface.CheckUserExistance(user.Email)

	if !existing {
		return errors.New("user already exists")
	}

	isvaild_email := uc.userVaildate.IsValidEmail(user.Email)
	ispassword_strong := uc.userVaildate.IsStrongPassword(user.Password)

	if !ispassword_strong || !isvaild_email {
		return errors.New("invalid password or email")
	}

	hashpassword := uc.userVaildate.Hashpassword(user.Password)
	user.Password = hashpassword

	err := uc.userinterface.Create("user_unverified",user)
	if err != nil {
		return err
	}

	err = uc.SendOTP(user.Email)

	if err != nil {
		return errors.New("failed to send OTP: " + err.Error())
	}

	return nil
}
func (uc *UserUsecase) SendOTP(email string) error {
	otp := uc.generateotp.GenerateRandomOTP()
	entry := domain.UserUnverified{
		Email:     email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := uc.generateotp.Send(email, otp); err != nil {
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
	_ = uc.userOTP.DeleteOTP(email)
	return true, nil
}
