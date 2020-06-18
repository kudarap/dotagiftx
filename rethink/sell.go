package rethink

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableSell = "sell"

// NewSell creates new instance of sell data store.
func NewSell(c *Client) core.SellStorage {
	if err := c.autoMigrate(tableSell); err != nil {
		log.Fatalf("could not create %s table: %s", tableSell, err)
	}

	return &sellStorage{c}
}

type sellStorage struct {
	db *Client
}

func (s *sellStorage) Find(o core.FindOpts) ([]core.Sell, error) {
	var res []core.Sell
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *sellStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{Filter: o.Filter, UserID: o.UserID}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *sellStorage) Get(id string) (*core.Sell, error) {
	row := &core.Sell{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.SellErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *sellStorage) Create(in *core.Sell) error {
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

func (s *sellStorage) Update(in *core.Sell) error {
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

func (s *sellStorage) table() r.Term {
	return r.Table(tableSell)
}
