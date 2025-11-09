package sql

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserInvalid       = errors.New("user invalid")
	ErrUserInternal      = errors.New("user internal error")

	ErrReputationGroupNotFound      = errors.New("reputation group not found")
	ErrReputationGroupAlreadyExists = errors.New("reputation group already exists")
	ErrReputationGroupInvalid       = errors.New("reputation group invalid")
	ErrReputationGroupInternal      = errors.New("reputation group internal error")

	ErrBalanceNotFound      = errors.New("balance not found")
	ErrBalanceAlreadyExists = errors.New("balance already exists")
	ErrBalanceNotEnough     = errors.New("balance not enough")
	ErrBalanceInternal      = errors.New("balance internal error")
	ErrBalanceInvalid       = errors.New("balance invalid")
)
