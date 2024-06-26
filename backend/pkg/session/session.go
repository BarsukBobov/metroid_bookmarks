package session

import (
	"crypto/rand"
	"encoding/hex"
	"metroid_bookmarks/internal/repository/sql/users"
)

const (
	AnonymousExpires     = 3600
	AuthenticatedExpires = 2592000
	CookieSessionName    = "X-Session"
	TokenLength          = 24
)

type Session struct {
	users.User
	Token   string `db:"token"   json:"token"`
	Expires int    `db:"expires" json:"expires"`
}

func (s *Session) IsAuthenticated() bool {
	return s.ID != 0
}

func (s *Session) IsAdmin() bool {
	return s.User.IsAdmin
}

func (s *Session) SetUser(user *users.User) {
	s.User = *user
	s.Expires = AuthenticatedExpires
}

func (s *Session) ResetUser() {
	s.User = users.User{}
	s.Expires = AnonymousExpires
}

func CreateToken() string {
	byteArray := make([]byte, TokenLength)
	if _, err := rand.Read(byteArray); err != nil {
		panic(err)
	}

	return hex.EncodeToString(byteArray)
}
