package authhasher

type AuthHasher interface {
	Compare(storedHash, authHash string) bool
	Hash(authHash string) (string, error)
}
