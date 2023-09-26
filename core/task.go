package core

import "time"

const (
	TaskKindVerifyDelivery  = 1
	TaskKindVerifyInventory = 2
)

type Task struct {
	ID        string     `json:"id"           db:"id,omitempty,indexed"`
	Kind      int        `json:"type"         db:"type,omitempty,indexed"`
	Value     string     `json:"value"        db:"value,omitempty"`
	CreatedAt *time.Time `json:"created_at"   db:"created_at,omitempty,indexed"`
	UpdatedAt *time.Time `json:"updated_at"   db:"updated_at,omitempty"`
}
