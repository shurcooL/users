package users

import (
	"html/template"
	"time"

	"golang.org/x/net/context"
)

// Service for users.
type Service interface {
	// GetAuthenticated fetches the currently authenticated user,
	// or nil if there is no authenticated user.
	GetAuthenticated(ctx context.Context) (*User, error)

	// Get fetches the specified user.
	Get(ctx context.Context, user UserSpec) (*User, error)

	// Edit the authenticated user.
	Edit(ctx context.Context, user *User) (*User, error)

	// CONSIDER: Delete user.
	//Delete(ctx context.Context, user UserSpec) error
}

type UserSpec struct {
	ID     uint64
	Domain string
}

// User represents a user.
type User struct {
	UserSpec
	Login     string
	AvatarURL template.URL
	HTMLURL   template.URL

	Name  string
	Email string // Public email.

	CreatedAt time.Time
	UpdatedAt time.Time

	SiteAdmin bool
}
