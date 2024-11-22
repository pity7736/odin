package usermodel

type User struct {
	email    string
	id       string
	password string
}

func New(email, id, password string) *User {
	return &User{email: email, id: id, password: password}
}

func (self *User) ID() string {
	return self.id
}

func (self *User) Email() string {
	return self.email
}

func (self *User) Password() string {
	return self.password
}

func (self *User) CheckPassword(password string) bool {
	return self.password == password
}
