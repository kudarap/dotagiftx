package rethink

import (
	"context"
	"errors"
	"log"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableInventory         = "inventory"
	inventoryFieldMarketID = "market_id"
)

var inventorySearchFields = []string{"id", "market_id"}

// NewInventory creates new instance of inventory data store.
func NewInventory(c *Client) *InventoryRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableInventory); err != nil {
		log.Fatalf("could not create %s table: %s", tableInventory, err)
	}

	if err := c.autoIndex(ctx, tableInventory, dotagiftx.Inventory{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableInventory, err)
	}

	return &InventoryRepository{c, inventorySearchFields}
}

type InventoryRepository struct {
	db            *Client
	keywordFields []string
}

func (s *InventoryRepository) Find(ctx context.Context, o dotagiftx.FindOpts) ([]dotagiftx.Inventory, error) {
	var res []dotagiftx.Inventory
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *InventoryRepository) Count(ctx context.Context, o dotagiftx.FindOpts) (num int, err error) {
	o = dotagiftx.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
	}
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	err = s.db.one(ctx, q.Count(), &num)
	return
}

// includeRelatedFields injects user details base on market foreign keys.
func (s *InventoryRepository) includeRelatedFields(q r.Term) r.Term {
	return q
	// return q.
	//	EqJoin(inventoryFieldMarketID, r.Table(tableMarket)).
	//	Map(func(t r.Term) r.Term {
	//		return t.Field("left").Merge(map[string]interface{}{
	//			tableMarket: t.Field("right"),
	//		})
	//	})
}

func (s *InventoryRepository) Get(ctx context.Context, id string) (*dotagiftx.Inventory, error) {
	row := &dotagiftx.Inventory{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.InventoryErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *InventoryRepository) GetByMarketID(ctx context.Context, marketID string) (*dotagiftx.Inventory, error) {
	var res []dotagiftx.Inventory
	var err error

	q := s.table().GetAllByIndex(inventoryFieldMarketID, marketID)
	if err = s.db.list(ctx, q, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, dotagiftx.InventoryErrNotFound
	}

	return &res[0], nil
}

func (s *InventoryRepository) Create(ctx context.Context, in *dotagiftx.Inventory) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *InventoryRepository) Update(ctx context.Context, in *dotagiftx.Inventory) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *InventoryRepository) table() r.Term {
	return r.Table(tableInventory)
}
