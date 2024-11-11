package user

type User struct {
	email string
}

func New(email string) *User {
	return &User{email: email}
}

func (u *User) Email() string {
	return u.email
}
