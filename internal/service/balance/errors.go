package balance

import "errors"

var (
	ErrBalanceNotFound  = errors.New("balance not found")
	ErrBalanceNotEnough = errors.New("balance not enough")
	ErrBalanceInternal  = errors.New("balance internal error")
	ErrBalanceInvalid   = errors.New("balance invalid")
)
