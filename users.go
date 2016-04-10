package users

import (
	"errors"
	"fmt"
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

// Static implementation of users.Service.
type Static struct{}

func (Static) Get(ctx context.Context, user UserSpec) (User, error) {
	const (
		ssg = "src.sourcegraph.com"
		gh  = "github.com"
		tw  = "twitter.com"
		ds  = "dmitri.shuralyov.com"
	)

	switch user {
	case UserSpec{ID: 678271, Domain: ssg}, UserSpec{ID: 1924134, Domain: gh}:
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{ /*{ID: 1, Domain: ds},*/ {ID: 1924134, Domain: gh}, {ID: 21361484, Domain: tw}},
			Login:     "shurcooL",
			Name:      "Dmitri Shuralyov",
			AvatarURL: "https://dmitri.shuralyov.com/avatar.jpg",
			HTMLURL:   "https://dmitri.shuralyov.com",
			SiteAdmin: true,
		}, nil
	case UserSpec{ID: 4332971, Domain: gh}: // Mee.
		return User{
			UserSpec:  user,
			Login:     "purple-snow",
			AvatarURL: "https://avatars.githubusercontent.com/u/4332971?v=3",
			HTMLURL:   "https://github.com/purple-snow",
		}, nil
	case UserSpec{ID: 43004, Domain: gh}: // pbakaus.
		return User{
			UserSpec:  user,
			Login:     "pbakaus",
			AvatarURL: "https://avatars.githubusercontent.com/u/43004?v=3",
			HTMLURL:   "https://github.com/pbakaus",
		}, nil
	case UserSpec{ID: 2, Domain: ds}: // Bernardo.
		return User{
			UserSpec:  user,
			Login:     "Bernardo",
			Name:      "Bernardo",
			AvatarURL: "https://secure.gravatar.com/avatar?d=mm&f=y&s=96",
		}, nil
	case UserSpec{ID: 3, Domain: ds}: // Michal Marcinkowski.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 15185890, Domain: tw}},
			Login:     "Michal Marcinkowski",
			Name:      "Michal Marcinkowski",
			AvatarURL: "https://secure.gravatar.com/avatar?d=mm&f=y&s=96",
		}, nil
	case UserSpec{ID: 4, Domain: ds}: // Anders Elfgren.
		return User{
			UserSpec:  user,
			Login:     "Anders Elfgren",
			Name:      "Anders Elfgren",
			AvatarURL: "https://secure.gravatar.com/avatar?d=mm&f=y&s=96",
		}, nil
	case UserSpec{ID: 5, Domain: ds}: // benp.
		return User{
			UserSpec:  user,
			Login:     "benp",
			AvatarURL: "https://secure.gravatar.com/avatar?d=mm&f=y&s=96",
		}, nil

	case UserSpec{ID: 678175, Domain: ssg}: // sqs.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 1976, Domain: gh}},
			Login:     "sqs",
			AvatarURL: "https://avatars.githubusercontent.com/u/1976?v=3",
			HTMLURL:   "https://github.com/sqs",
		}, nil
	case UserSpec{ID: 678177, Domain: ssg}: // slimsag.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 3173176, Domain: gh}},
			Login:     "slimsag",
			AvatarURL: "https://avatars.githubusercontent.com/u/3173176?v=3",
			HTMLURL:   "https://github.com/slimsag",
		}, nil
	case UserSpec{ID: 678180, Domain: ssg}: // keegancsmith.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 187831, Domain: gh}},
			Login:     "keegancsmith",
			AvatarURL: "https://avatars.githubusercontent.com/u/187831?v=3",
			HTMLURL:   "https://github.com/keegancsmith",
		}, nil
	case UserSpec{ID: 678179, Domain: ssg}: // renfredxh.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 3800339, Domain: gh}},
			Login:     "renfredxh",
			AvatarURL: "https://avatars.githubusercontent.com/u/3800339?v=3",
			HTMLURL:   "https://github.com/renfredxh",
		}, nil
	case UserSpec{ID: 678176, Domain: ssg}: // nicot.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 3722365, Domain: gh}},
			Login:     "nicot",
			AvatarURL: "https://avatars.githubusercontent.com/u/3722365?v=3",
			HTMLURL:   "https://github.com/nicot",
		}, nil
	case UserSpec{ID: 678357, Domain: ssg}: // rothfels.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 1095573, Domain: gh}},
			Login:     "rothfels",
			AvatarURL: "https://avatars.githubusercontent.com/u/1095573?v=3",
			HTMLURL:   "https://github.com/rothfels",
		}, nil
	case UserSpec{ID: 678225, Domain: ssg}: // beyang.
		return User{
			UserSpec:  user,
			Elsewhere: []UserSpec{{ID: 1646931, Domain: gh}},
			Login:     "beyang",
			AvatarURL: "https://avatars.githubusercontent.com/u/1646931?v=3",
			HTMLURL:   "https://github.com/beyang",
		}, nil
	default:
		return User{}, fmt.Errorf("user %v not found", user)
	}
}

func (s Static) GetAuthenticated(ctx context.Context) (*UserSpec, error) {
	// TEMP, HACK: Pretend I'm logged in (for testing).
	//return &UserSpec{ID: 1924134, Domain: "github.com"}, nil

	// Authenticated user not yet supported.
	return nil, nil
}

func (Static) Edit(ctx context.Context, er EditRequest) (User, error) {
	return User{}, errors.New("Edit is not implemented")
}
