package users

import (
	"html/template"

	"golang.org/x/net/context"
)

// Service for users.
type Service interface {
	// CONSIDER: AuthenticatedUser()
	CurrentUser(ctx context.Context) (*User, error)

	// Get fetches a user.
	// CONSIDER: Passing the empty string will fetch the authenticated user.
	Get(ctx context.Context, u UserSpec) (*User, error)

	// Edit the authenticated user.
	Edit(user *User) (*User, error)

	// CONSIDER: Delete user.
	//Delete(user UserSpec) error
}

// User represents a user.
type User struct {
	UserSpec
	Login     string
	AvatarURL template.URL
	HTMLURL   template.URL

	Name  string
	Email string // Public email.

	CreatedAt
	UpdatedAt

	SiteAdmin bool
}

type UserSpec struct {
	ID     uint64
	Domain string
}
