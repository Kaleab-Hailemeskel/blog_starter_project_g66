package domain
type UserRepository interface {
	Create(user *User) error
	Login(email, password string) (*User, error) //authenticate and return user
	FetchByEmail(email string) (bool, error) //checks if user exisits or not
}