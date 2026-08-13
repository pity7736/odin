package sessionmodel

import (
	"time"

	"raiseexception.dev/odin/src/shared/utils"
)

const DefaultTTL = 30 * 24 * time.Hour
const tokenLength uint8 = 50

type Session struct {
	expiresAt time.Time
	createdAt time.Time
	token     string
	userID    string
}

func New(userID string, ttl time.Duration) (*Session, error) {
	token, err := utils.RandomString(tokenLength)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Session{
		token:     token,
		userID:    userID,
		createdAt: now,
		expiresAt: now.Add(ttl),
	}, nil
}

func NewFromRepository(token, userID string, createdAt, expiresAt time.Time) *Session {
	return &Session{
		token:     token,
		userID:    userID,
		createdAt: createdAt,
		expiresAt: expiresAt,
	}
}

func (self *Session) Token() string {
	return self.token
}

func (self *Session) UserID() string {
	return self.userID
}

func (self *Session) CreatedAt() time.Time {
	return self.createdAt
}

func (self *Session) ExpiresAt() time.Time {
	return self.expiresAt
}

func (self *Session) IsExpired() bool {
	return time.Now().After(self.expiresAt)
}

func (self *Session) Extend(ttl time.Duration) {
	self.expiresAt = time.Now().Add(ttl)
}
