// Package fs implements an in-memory users.Store backed by a virtual filesystem.
package fs

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/shurcooL/users"
	"golang.org/x/net/webdav"
)

// NewStore creates an in-memory users.Store, backed by virtual filesystem
// using root for storage.
func NewStore(root webdav.FileSystem) (users.Store, error) {
	s := &store{
		fs:    root,
		users: make(map[users.UserSpec]users.User),
	}
	err := s.load()
	if err != nil {
		return nil, err
	}
	return s, nil
}

type store struct {
	mu    sync.Mutex
	fs    webdav.FileSystem
	users map[users.UserSpec]users.User
}

func (s *store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := s.fs.OpenFile(context.Background(), "users", os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for {
		var u user
		err := dec.Decode(&u)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		user := u.User()
		s.users[user.UserSpec] = user
	}
	return nil
}

func (s *store) Create(ctx context.Context, user users.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user.UserSpec]; ok {
		return os.ErrExist
	}

	// Commit to storage first, returning error on failure.
	f, err := s.fs.OpenFile(ctx, "users", os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(fromUser(user))
	if err != nil {
		return err
	}

	// Commit to memory second.
	s.users[user.UserSpec] = user
	return nil
}

func (s *store) Get(_ context.Context, user users.UserSpec) (users.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[user]
	if !ok {
		return users.User{}, os.ErrNotExist
	}
	return u, nil
}
