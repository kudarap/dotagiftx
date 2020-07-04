package rethink

import (
	"log"

	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTrack       = "track"
	trackFieldItemID = "item_id"
)

// NewTrack creates new instance of track data store.
func NewTrack(c *Client) *trackStorage {
	if err := c.autoMigrate(tableTrack); err != nil {
		log.Fatalf("could not create %s table: %s", tableTrack, err)
	}

	if err := c.createIndex(tableMarket, trackFieldItemID); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &trackStorage{c, []string{"item_id"}}
}

type trackStorage struct {
	db            *Client
	keywordFields []string
}

func (s *trackStorage) Find(o core.FindOpts) ([]core.Track, error) {
	var res []core.Track
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *trackStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
	}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *trackStorage) Get(id string) (*core.Track, error) {
	row := &core.Track{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.TrackErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *trackStorage) Create(in *core.Track) error {
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

func (s *trackStorage) table() r.Term {
	return r.Table(tableTrack)
}
