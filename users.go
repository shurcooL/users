// Package users provides a users service definition.
package users

import (
	"html/template"
	"time"

	"golang.org/x/net/context"
)

// Service for users.
type Service interface {
	// Get fetches the specified user.
	Get(ctx context.Context, user UserSpec) (User, error)

	// TODO: Consider zero value UserSpec instead of nil for no user.
	// GetAuthenticated fetches the currently authenticated user specification,
	// or nil if there is no authenticated user.
	GetAuthenticated(ctx context.Context) (*UserSpec, error)

	// Edit the authenticated user.
	Edit(ctx context.Context, er EditRequest) (User, error)

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
	Elsewhere []UserSpec // THINK: Consider merging Elsewhere with root-most UserSpec. Maybe even use a set, order of linked accounts shouldn't matter, should it?

	Login     string
	Name      string
	Email     string // Public email.
	AvatarURL template.URL
	HTMLURL   template.URL

	CreatedAt time.Time
	UpdatedAt time.Time

	SiteAdmin bool
}

// EditRequest represents a request to edit a user.
type EditRequest struct {
	// Currently nothing, but editable fields will be added here in the future.
}
