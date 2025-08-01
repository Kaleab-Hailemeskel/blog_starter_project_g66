package domain


type IUserRepository interface {
	Create(username string, user *User) error
	Login(email, password string) (*User, error) //authenticate and return user
	FetchByEmail(email string) (bool, error)     //checks if user exisits or not
}

type IUserValidation interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) bool
	Hashpassword(password string) string
	ComparePassword(userPassword, password string) error
}
type IUserOTP interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
}