package user

type User struct {
	email string
}

func New(email string) *User {
	return &User{email: email}
}
