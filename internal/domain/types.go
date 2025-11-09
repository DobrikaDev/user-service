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

func SexToSex(sex string) Sex {
	switch sex {
	case "male":
		return SexMale
	case "female":
		return SexFemale
	default:
		return SexUnknown
	}
}
func RoleToUserRole(role string) UserRole {
	switch role {
	case "pending":
		return UserRoleAdmin
	case "user":
		return UserRoleUser
	default:
		return UserRoleUser
	}
}

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
