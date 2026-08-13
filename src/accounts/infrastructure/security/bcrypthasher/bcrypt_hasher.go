package bcrypthasher

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func New() BcryptHasher {
	return BcryptHasher{}
}

func (self BcryptHasher) Compare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
