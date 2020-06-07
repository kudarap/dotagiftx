package rethink

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableItem = "item"

// NewItem creates new instance of item data store.
func NewItem(c *Client) core.ItemStorage {
	if err := c.autoMigrate(tableItem); err != nil {
		log.Fatalf("could not create %s table: %s", tableItem, err)
	}

	return &itemStorage{c}
}

type itemStorage struct {
	db *Client
}

func (s *itemStorage) Find(o core.FindOpts) ([]core.Item, error) {
	var res []core.Item
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *itemStorage) Get(id string) (*core.Item, error) {
	row := &core.Item{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.ItemErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *itemStorage) Create(in *core.Item) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *itemStorage) Update(in *core.Item) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return errors.New(core.StorageMergeErr, err)
	}

	return nil
}

func (s *itemStorage) table() r.Term {
	return r.Table(tableItem)
}
