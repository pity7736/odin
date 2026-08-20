package bcrypthasher

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func New() BcryptHasher {
	return BcryptHasher{}
}

func (self BcryptHasher) Compare(storedHash, authHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(authHash)) == nil
}

func (self BcryptHasher) Hash(authHash string) (string, error) {
	digest, err := bcrypt.GenerateFromPassword([]byte(authHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(digest), nil
}
