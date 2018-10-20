package auth

import (
	"github.com/boreq/blogs/database"
)

type User interface {
	// IsAuthenticated returns true if the user has logged in.
	IsAuthenticated() bool

	// GetUser returns an associated database entry or nil if the user is
	// not authenticated.
	GetUser() *database.User
}

// authenticatedUser is used to represent a logged in user.
type authenticatedUser struct {
	user database.User
}

func (u authenticatedUser) IsAuthenticated() bool {
	return true
}

func (u authenticatedUser) GetUser() *database.User {
	return &u.user
}

// anonymousUser is used to represent a user that is not logged in.
type anonymousUser struct{}

func (u anonymousUser) IsAuthenticated() bool {
	return false
}

func (u anonymousUser) GetUser() *database.User {
	return nil
}
