package reputationgroup

import "errors"

var (
	ErrReputationGroupNotFound      = errors.New("reputation group not found")
	ErrReputationGroupAlreadyExists = errors.New("reputation group already exists")
	ErrReputationGroupInvalid       = errors.New("reputation group invalid")
	ErrReputationGroupInternal      = errors.New("reputation group internal error")
)
