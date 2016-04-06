package users_test

import "github.com/shurcooL/users"

// Test that users.Static{} implements users.Service.
var _ users.Service = users.Static{}
