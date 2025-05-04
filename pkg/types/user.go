package types

import "time"

type User struct {
	ID             int        `json:"id,omitempty"`
	Tokens         int        `json:"tokens"`
	TokensCap      int        `json:"tokens_cap"`
	RatePerMinute  int        `json:"rate"`
	CreatedAt      *time.Time `json:"-"`
	UpdatedAt      *time.Time `json:"-"`
	LastAdditionAt *time.Time `json:"-"`
}
