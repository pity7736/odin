package authhasher

type AuthHasher interface {
	Compare(storedHash, authHash string) bool
}
