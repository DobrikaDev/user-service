package user

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserInternal = errors.New("user internal error")
var ErrUserInvalid = errors.New("user invalid")
