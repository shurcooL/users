// Package asanaapi implements users.Service using Asana API client.
package asanaapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/shurcooL/users"
	"github.com/tambet/go-asana/asana"
)

// NewService creates a Asana-backed users.Service using given Asana client.
// At this time it infers the current user from the client (its authentication info),
// and cannot be used to serve multiple users.
func NewService(client *asana.Client) (users.Service, error) {
	if client == nil {
		client = asana.NewClient(nil)
	}
	var currentUser users.User
	switch u, err := client.GetAuthenticatedUser(context.Background(), nil); err {
	case nil:
		currentUser = asanaUser(u)
	case asana.ErrUnauthorized:
		// There's no authenticated user.
		currentUser = users.User{}
	default:
		return nil, err
	}
	return service{
		cl:          client,
		currentUser: currentUser,
	}, nil
}

type service struct {
	cl *asana.Client

	currentUser users.User
}

func (s service) Get(ctx context.Context, user users.UserSpec) (users.User, error) {
	if user.Domain != "app.asana.com" {
		return users.User{}, fmt.Errorf("user %v not found", user)
	}

	u, err := s.cl.GetUserByID(ctx, int64(user.ID), nil)
	if err != nil {
		return users.User{}, err
	}
	return asanaUser(u), nil
}

func (s service) GetAuthenticated(ctx context.Context) (users.User, error) {
	return s.currentUser, nil
}

func (s service) GetAuthenticatedSpec(ctx context.Context) (users.UserSpec, error) {
	return s.currentUser.UserSpec, nil
}

func (s service) Edit(ctx context.Context, er users.EditRequest) (users.User, error) {
	return users.User{}, errors.New("Edit is not implemented")
}

func asanaUser(user asana.User) users.User {
	return users.User{
		UserSpec: users.UserSpec{
			ID:     uint64(user.ID),
			Domain: "app.asana.com",
		},
		Login:     user.Name,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.Photo["image_128x128"],
	}
}
