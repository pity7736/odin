package repositories

type SessionRepository interface {
	Add(token string) error
}
