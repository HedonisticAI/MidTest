package postgres_repo

import "errors"

var ErrBadAuth = errors.New("bad login or password")
