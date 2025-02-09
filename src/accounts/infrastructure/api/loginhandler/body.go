package loginhandler

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
