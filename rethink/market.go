package rethink

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableMarket = "market"

// NewMarket creates new instance of market data store.
func NewMarket(c *Client) core.MarketStorage {
	if err := c.autoMigrate(tableMarket); err != nil {
		log.Fatalf("could not create %s table: %s", tableMarket, err)
	}

	return &marketStorage{c, []string{"name", "hero", "origin"}}
}

type marketStorage struct {
	db            *Client
	keywordFields []string
}

func (s *marketStorage) Find(o core.FindOpts) ([]core.Market, error) {
	var res []core.Market
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *marketStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
	}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *marketStorage) Get(id string) (*core.Market, error) {
	row := &core.Market{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.MarketErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *marketStorage) Create(in *core.Market) error {
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

func (s *marketStorage) Update(in *core.Market) error {
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

func (s *marketStorage) table() r.Term {
	return r.Table(tableMarket)
}
