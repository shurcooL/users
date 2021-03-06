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
// At this time it infers the current user from the client (its authentication info),
// and cannot be used to serve multiple users.
func NewService(client *github.Client) (users.Service, error) {
	if client == nil {
		client = github.NewClient(nil)
	}
	var currentUser users.User
	switch u, _, err := client.Users.Get(context.Background(), ""); {
	case err == nil:
		currentUser = ghUser(u)
	case isUnauthorized(err):
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

// isUnauthorized reports whether err is an unauthorized error response from GitHub.
func isUnauthorized(err error) bool {
	e, ok := err.(*github.ErrorResponse)
	return ok && e.Response.StatusCode == http.StatusUnauthorized
}

type service struct {
	cl *github.Client

	currentUser users.User
}

func (s service) Get(ctx context.Context, user users.UserSpec) (users.User, error) {
	if user.Domain != "github.com" {
		return users.User{}, fmt.Errorf("user %v not found", user)
	}

	ghUser, _, err := s.cl.Users.GetByID(ctx, int64(user.ID))
	if err != nil {
		return users.User{}, err
	}
	if ghUser.ID == nil || ghUser.Login == nil || ghUser.AvatarURL == nil || ghUser.HTMLURL == nil {
		return users.User{}, fmt.Errorf("github user missing fields: %#v", ghUser)
	}
	if uint64(*ghUser.ID) != user.ID {
		return users.User{}, fmt.Errorf("got github user ID %v, but wanted ID %v", *ghUser.ID, user.ID)
	}
	return users.User{
		UserSpec:  user,
		Login:     *ghUser.Login,
		AvatarURL: *ghUser.AvatarURL,
		HTMLURL:   *ghUser.HTMLURL,
	}, nil
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
