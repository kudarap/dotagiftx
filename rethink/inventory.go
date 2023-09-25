package rethink

import (
	"log"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableInventory         = "inventory"
	inventoryFieldMarketID = "market_id"
)

var inventorySearchFields = []string{"id", "market_id"}

// NewInventory creates new instance of inventory data store.
func NewInventory(c *Client) core.InventoryStorage {
	if err := c.autoMigrate(tableInventory); err != nil {
		log.Fatalf("could not create %s table: %s", tableInventory, err)
	}

	if err := c.autoIndex(tableInventory, core.Inventory{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableInventory, err)
	}

	return &inventoryStorage{c, inventorySearchFields}
}

type inventoryStorage struct {
	db            *Client
	keywordFields []string
}

func (s *inventoryStorage) Find(o core.FindOpts) ([]core.Inventory, error) {
	var res []core.Inventory
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *inventoryStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
	}
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	err = s.db.one(q.Count(), &num)
	return
}

// includeRelatedFields injects user details base on market foreign keys.
func (s *inventoryStorage) includeRelatedFields(q r.Term) r.Term {
	return q
	//return q.
	//	EqJoin(inventoryFieldMarketID, r.Table(tableMarket)).
	//	Map(func(t r.Term) r.Term {
	//		return t.Field("left").Merge(map[string]interface{}{
	//			tableMarket: t.Field("right"),
	//		})
	//	})
}

func (s *inventoryStorage) Get(id string) (*core.Inventory, error) {
	row := &core.Inventory{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.InventoryErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *inventoryStorage) GetByMarketID(marketID string) (*core.Inventory, error) {
	var res []core.Inventory
	var err error

	q := s.table().GetAllByIndex(inventoryFieldMarketID, marketID)
	if err = s.db.list(q, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, core.InventoryErrNotFound
	}

	return &res[0], nil
}

func (s *inventoryStorage) Create(in *core.Inventory) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *inventoryStorage) Update(in *core.Inventory) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return errors.New(core.StorageMergeErr, err)
	}

	return nil
}

func (s *inventoryStorage) table() r.Term {
	return r.Table(tableInventory)
}
