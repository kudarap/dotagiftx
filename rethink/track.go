package rethink

import (
	"log"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTrack          = "track"
	trackFieldItemID    = "item_id"
	trackFieldType      = "type"
	trackFieldCreatedAt = "created_at"
)

// NewTrack creates new instance of track data store.
func NewTrack(c *Client) *trackStorage {
	if err := c.autoMigrate(tableTrack); err != nil {
		log.Fatalf("could not create %s table: %s", tableTrack, err)
	}

	if err := c.autoIndex(tableTrack, dotagiftx.Track{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &trackStorage{c, []string{"item_id"}}
}

type trackStorage struct {
	db            *Client
	keywordFields []string
}

func (s *trackStorage) Find(o dotagiftx.FindOpts) ([]dotagiftx.Track, error) {
	var res []dotagiftx.Track
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *trackStorage) Count(o dotagiftx.FindOpts) (num int, err error) {
	o = dotagiftx.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
	}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *trackStorage) Get(id string) (*dotagiftx.Track, error) {
	row := &dotagiftx.Track{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dotagiftx.TrackErrNotFound
		}

		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *trackStorage) Create(in *dotagiftx.Track) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

const last7days = 604800

// TopKeywords returns top recent searched keywords.
//
/*
	var thisWeek = 604800;
	r.db('d2g_production').table('track')
	 .between(r.now().sub(thisWeek), r.now(), {index: 'created_at'})
	 .filter({ type: 's' })
	 .group('keyword')
	 .count()
	 .ungroup()
	 .orderBy(r.desc('reduction'))
	 .limit(12)
	 .map(function(doc) {
	   return {
		 keyword: doc('group'),
		 score: doc('reduction'),
	   }
	 })
*/
func (s *trackStorage) TopKeywords() ([]dotagiftx.SearchKeywordScore, error) {
	now := r.Now()
	q := s.table().Between(now.Sub(last7days), now, r.BetweenOpts{Index: trackFieldCreatedAt}).
		Filter(map[string]interface{}{"type": "s"}).
		Group(r.Row.Field("keyword").Downcase()).
		Count().
		Ungroup().
		OrderBy(r.Desc("reduction")).
		Limit(12).
		Map(func(doc r.Term) interface{} {
			return map[string]interface{}{
				"Keyword": doc.Field("group"),
				"Score":   doc.Field("reduction"),
			}
		})

	var res []dotagiftx.SearchKeywordScore
	if err := s.db.list(q, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *trackStorage) table() r.Term {
	return r.Table(tableTrack)
}
