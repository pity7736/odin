package bcrypthasher

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func New() BcryptHasher {
	return BcryptHasher{}
}

func (self BcryptHasher) Compare(storedHash, authHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(authHash)) == nil
}
