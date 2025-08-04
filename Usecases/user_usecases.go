package usecases

import domain "blog_starter_project_g66/Domain"

type UserUseCase struct {
	UserDataBase domain.IUserRepository
}

func NewUserUseCase(blogDB domain.IBlogRepository, userDB domain.IUserRepository) *UserUseCase {
	return &UserUseCase{
		UserDataBase: userDB,
	}
}
func (uuc *UserUseCase) DemoteUser(userEmail string) error {
	return uuc.UserDataBase.DemoteUser(userEmail)
}
func (uuc *UserUseCase) PromoteUser(userEmail string) error {
	return uuc.UserDataBase.PromoteUser(userEmail)
}
