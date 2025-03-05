package domain

type User struct {
	ID int
	Username string
	Email string
	Password string
}

func NewUser(id int,username string, email string, password string) *User{
	return &User{
		ID:id,
		Username:username,
		Email:email,
		Password: password,
	}
}

func (u *User) GetEmail() string {
	return u.Email
}
