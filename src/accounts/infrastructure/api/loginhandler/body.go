package loginhandler

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
