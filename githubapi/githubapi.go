// Package githubapi implements users.Service using GitHub API client.
package githubapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/shurcooL/users"
)

// NewService creates a GitHub-backed users.Service using given GitHub client.
// At this time it infers the current user from the client (its authentication info), and cannot be used to serve multiple users.
func NewService(client *github.Client) users.Service {
	if client == nil {
		client = github.NewClient(nil)
	}

	s := service{
		cl: client,
	}

	if u, _, err := s.cl.Users.Get(""); err == nil {
		s.currentUser = ghUser(u)
		s.currentUserErr = nil
	} else if ghErr, ok := err.(*github.ErrorResponse); ok && ghErr.Response.StatusCode == http.StatusUnauthorized {
		// There's no authenticated user.
		s.currentUser = users.User{}
		s.currentUserErr = nil
	} else {
		s.currentUser = users.User{}
		s.currentUserErr = err
	}

	return s
}

type service struct {
	cl *github.Client

	currentUser    users.User
	currentUserErr error
}

func (s service) Get(ctx context.Context, user users.UserSpec) (users.User, error) {
	if user.Domain != "github.com" {
		return users.User{}, fmt.Errorf("user %v not found", user)
	}

	ghUser, _, err := s.cl.Users.GetByID(int(user.ID))
	if err != nil {
		return users.User{}, err
	}
	if ghUser.Login == nil || ghUser.AvatarURL == nil || ghUser.HTMLURL == nil {
		return users.User{}, fmt.Errorf("github user missing fields: %#v", ghUser)
	}
	return users.User{
		UserSpec:  user,
		Login:     *ghUser.Login,
		AvatarURL: *ghUser.AvatarURL,
		HTMLURL:   *ghUser.HTMLURL,
	}, nil
}

func (s service) GetAuthenticated(ctx context.Context) (users.User, error) {
	return s.currentUser, s.currentUserErr
}

func (s service) GetAuthenticatedSpec(ctx context.Context) (users.UserSpec, error) {
	return s.currentUser.UserSpec, s.currentUserErr
}

func (s service) Edit(ctx context.Context, er users.EditRequest) (users.User, error) {
	return users.User{}, errors.New("Edit is not implemented")
}

func ghUser(user *github.User) users.User {
	return users.User{
		UserSpec: users.UserSpec{
			ID:     uint64(*user.ID),
			Domain: "github.com",
		},
		Login:     *user.Login,
		AvatarURL: *user.AvatarURL,
		HTMLURL:   *user.HTMLURL,
	}
}
