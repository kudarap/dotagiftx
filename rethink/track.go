package rethink

import (
	"context"
	"errors"
	"log"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTrack          = "track"
	trackFieldItemID    = "item_id"
	trackFieldType      = "type"
	trackFieldCreatedAt = "created_at"
)

// NewTrack creates new instance of track data store.
func NewTrack(c *Client) *TrackRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableTrack); err != nil {
		log.Fatalf("could not create %s table: %s", tableTrack, err)
	}

	if err := c.autoIndex(ctx, tableTrack, dotagiftx2.Track{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &TrackRepository{c, []string{"item_id"}}
}

type TrackRepository struct {
	db            *Client
	keywordFields []string
}

func (s *TrackRepository) Find(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Track, error) {
	var res []dotagiftx2.Track
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *TrackRepository) Count(ctx context.Context, o dotagiftx2.FindOpts) (num int, err error) {
	o = dotagiftx2.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
	}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(ctx, q.Count(), &num)
	return
}

func (s *TrackRepository) Get(ctx context.Context, id string) (*dotagiftx2.Track, error) {
	row := &dotagiftx2.Track{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.TrackErrNotFound
		}

		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *TrackRepository) Create(ctx context.Context, in *dotagiftx2.Track) error {
	_, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}
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
func (s *TrackRepository) TopKeywords(ctx context.Context) ([]dotagiftx2.SearchKeywordScore, error) {
	now := r.Now()
	q := s.table().Between(now.Sub(last7days), now, r.BetweenOpts{Index: trackFieldCreatedAt}).
		Filter(map[string]any{"type": "s"}).
		Group(r.Row.Field("keyword").Downcase()).
		Count().
		Ungroup().
		OrderBy(r.Desc("reduction")).
		Limit(12).
		Map(func(doc r.Term) any {
			return map[string]any{
				"Keyword": doc.Field("group"),
				"Score":   doc.Field("reduction"),
			}
		})

	var res []dotagiftx2.SearchKeywordScore
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TrackRepository) table() r.Term {
	return r.Table(tableTrack)
}
