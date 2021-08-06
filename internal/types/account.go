package types

type AccountData struct {
	ID        uint64 `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}
