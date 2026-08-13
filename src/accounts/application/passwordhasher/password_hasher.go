package passwordhasher

type PasswordHasher interface {
	Compare(hashedPassword, password string) bool
}
