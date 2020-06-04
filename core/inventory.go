package core

import (
	"context"
	"time"
)

const (
	InventoryStatusPending  InventoryStatus = 100
	InventoryStatusReserved InventoryStatus = 200
	InventoryStatusSold     InventoryStatus = 300
	InventoryStatusDeleted  InventoryStatus = 400
)

type (
	InventoryStatus uint

	Inventory struct {
		ID        string          `json:"id"         db:"id,omitempty"`
		Item      string          `json:"item"       db:"item,omitempty"`
		Hero      string          `json:"hero"       db:"hero,omitempty"`
		Price     float64         `json:"price"      db:"price,omitempty"`
		Currency  string          `json:"currency"   db:"currency,omitempty"`
		Notes     string          `json:"notes"      db:"notes,omitempty"`
		Status    InventoryStatus `json:"status"     db:"status,omitempty"`
		CreatedAt *time.Time      `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time      `json:"updated_at" db:"updated_at,omitempty"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"-"`
	}

	InventoryService interface {
		Inventories(opts FindOpts) ([]Inventory, error)
		Inventory(id string) (*Inventory, error)
		Create(context.Context, *Inventory) error
		Update(context.Context, *Inventory) error
		Delete(ctx context.Context, id string) error
	}

	InventoryStorage interface {
		Find(opts FindOpts) ([]Inventory, error)
		Get(id string) (*Inventory, error)
		Create(*Inventory) error
		Update(*Inventory) error
	}
)
