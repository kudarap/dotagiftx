package rethink

import (
	"fmt"
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableItem     = "item"
	itemFieldName = "name"
	itemFieldSlug = "slug"
)

var itemSearchFields = []string{"name", "hero", "origin", "rarity"}

// NewItem creates new instance of item data store.
func NewItem(c *Client) core.ItemStorage {
	if err := c.autoMigrate(tableItem); err != nil {
		log.Fatalf("could not create %s table: %s", tableItem, err)
	}

	if err := c.createIndex(tableItem, itemFieldSlug); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableItem, err)
	}

	return &itemStorage{c, itemSearchFields}
}

type itemStorage struct {
	db            *Client
	keywordFields []string
}

func (s *itemStorage) Find(o core.FindOpts) ([]core.Item, error) {
	var res []core.Item
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *itemStorage) Count(o core.FindOpts) (num int, err error) {
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

func (s *itemStorage) Get(id string) (*core.Item, error) {
	row, _ := s.GetBySlug(id)
	if row != nil {
		return row, nil
	}

	row = &core.Item{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.ItemErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *itemStorage) GetBySlug(slug string) (*core.Item, error) {
	row := &core.Item{}
	q := s.table().GetAllByIndex(itemFieldSlug, slug)
	if err := s.db.one(q, row); err != nil {
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
	in.ID = ""
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

func (s *itemStorage) IsItemExist(name string) error {
	/*
		r.table('item').filter(function(doc) {
		  return doc.getField('name').match('(?i)^Gothic')
		})
	*/
	q := s.table().Filter(func(t r.Term) r.Term {
		// Matches exact name and non case sensitive.
		return t.Field(itemFieldName).Match(fmt.Sprintf("(?i)^%s$", name))
	})
	var n int
	if err := s.db.one(q.Count(), &n); err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	if n != 0 {
		return core.ItemErrCreateItemExists
	}

	return nil
}

func (s *itemStorage) AddViewCount(id string) error {
	cur, err := s.Get(id)
	if err != nil {
		return err
	}

	cur.ViewCount++
	if err := s.Update(cur); err != nil {
		return err
	}

	if err := s.updateCatalogViewCount(id, cur.ViewCount); err != nil {
		return err
	}

	return nil
}

func (s *itemStorage) updateCatalogViewCount(itemID string, viewCount int) error {
	q := r.Table(tableCatalog).Get(itemID).Update(&core.Catalog{ViewCount: viewCount})
	return s.db.update(q)
}

func (s *itemStorage) table() r.Term {
	return r.Table(tableItem)
}
