package domain

type IUserRepository interface {
	Create(user *User) error
	Login(email, password string) (*User, error) //authenticate and return user
	FetchByEmail(email string) (bool, error)     //checks if user exisits or not
}

type IUserValidation interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) bool
	Hashpassword(password string) string
	ComparePassword(userPassword, password string) error
}
