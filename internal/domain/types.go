package domain

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type Sex string

const (
	SexMale    Sex = "male"
	SexFemale  Sex = "female"
	SexUnknown Sex = "unknown"
)

func (s UserRole) String() string {
	return string(s)
}

func (s Sex) String() string {
	return string(s)
}

type BalanceOperationType string

const (
	BalanceOperationTypeDeposit  BalanceOperationType = "deposit"
	BalanceOperationTypeWithdraw BalanceOperationType = "withdraw"
)

func (t BalanceOperationType) String() string {
	return string(t)
}
