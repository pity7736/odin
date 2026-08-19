package loginhandler

type LoginBody struct {
	Email    string `json:"email"`
	AuthHash string `json:"auth_hash"`
}
