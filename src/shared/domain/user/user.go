package user

type User struct {
	email string
	id    string
}

func New(email, id string) *User {
	return &User{email: email, id: id}
}

func (self *User) ID() string {
	return self.id
}

func (self *User) Email() string {
	return self.email
}
