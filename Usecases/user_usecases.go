package usecases

import (
	"blog_starter_project_g66/Domain"
	"errors"
)

type UserUsecase struct{
	userinterface domain.IUserRepository
	userVaildate domain.IUserValidation
}

func NewUserUsecase(ui domain.IUserRepository) *UserUsecase{
	return &UserUsecase{
		userinterface: ui,
	}
}

func (uc *UserUsecase)HandleRegistration(user *domain.User)error{
	existing, _ := uc.userinterface.FetchByEmail(user.Email)

	if !existing{
		return errors.New("user already exists")
	}

	isvaild_email := uc.userVaildate.IsValidEmail(user.Email)
	ispassword_strong := uc.userVaildate.IsStrongPassword(user.Password)

	if !ispassword_strong || !isvaild_email {
		return errors.New("invalid password or email")
	}

	hashpassword := uc.userVaildate.Hashpassword(user.Password)
	user.Password = hashpassword

	err := uc.userinterface.Create(user)
	return err
}
