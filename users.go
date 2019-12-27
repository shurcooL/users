// Package users provides a users service definition.
package users

import (
	"context"
)

// Service for users.
type Service interface {
	// Get fetches the specified user.
	Get(ctx context.Context, user UserSpec) (User, error)

	// GetAuthenticatedSpec fetches the currently authenticated
	// user specification, or UserSpec{ID: 0} if there is no
	// authenticated user. A non-nil error is returned if the
	// the authentication process was not able to successfully
	// determine if a user or no user is currently authenticated.
	GetAuthenticatedSpec(ctx context.Context) (UserSpec, error)

	// GetAuthenticated fetches the currently authenticated user,
	// or User{UserSpec: UserSpec{ID: 0}} if there is no authenticated user.
	GetAuthenticated(ctx context.Context) (User, error)

	// Edit the authenticated user.
	Edit(ctx context.Context, er EditRequest) (User, error)
}

// Store for users.
type Store interface {
	// Create creates the specified user.
	// UserSpec must specify a valid (i.e., non-zero) user.
	// It returns os.ErrExist if the user already exists.
	Create(ctx context.Context, user User) error

	// Get fetches the specified user.
	Get(ctx context.Context, user UserSpec) (User, error)
}

// UserSpec specifies a user.
// ID value 0 represents no user. Valid users may not use 0 as their ID.
type UserSpec struct {
	ID     uint64
	Domain string
}

// THINK: (2016-03-06) Consider merging Elsewhere with root-most UserSpec.
//        Maybe even use a set, order of linked accounts shouldn't matter, should it?
//
//        (2020-01-12) I'm not ready to commit to supporting multiple linked identities,
//        so for now I'm documenting that Elsewhere is metadata only and not used for
//        user equivalence or authentication purposes. But will keep thinking about this.

// User represents a user.
type User struct {
	// UserSpec is the primary user identity. It is mostly used for user equivalence,
	// and potentially for authentication for well-known domains (e.g., "github.com").
	UserSpec
	// CanonicalMe is an optional canonical user profile URL. When a non-zero value,
	// it is used for identifying users that authenticate via the IndieAuth protocol.
	CanonicalMe string

	// Elsewhere represents alternative user identities. This information is not used for
	// user equivalence or authentication purposes, but can be used for display purposes.
	Elsewhere []UserSpec

	Login     string
	Name      string
	Email     string // Public email.
	AvatarURL string
	HTMLURL   string

	SiteAdmin bool
}

// EditRequest represents a request to edit a user.
type EditRequest struct {
	// Currently nothing, but editable fields will be added here in the future.
}
