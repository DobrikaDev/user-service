package domain

type ReputationGroup struct {
	ID             int     `json:"id" db:"id"`
	Name           string  `json:"name" db:"name"`
	Description    string  `json:"description" db:"description"`
	Coefficient    float64 `json:"coefficient" db:"coefficient"`
	ReputationNeed int     `json:"reputation_need" db:"reputation_need"`
}
