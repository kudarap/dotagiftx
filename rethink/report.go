package rethink

import (
	"context"
	"errors"
	"log"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableReport       = "report"
	reportFieldUserID = "user_id"
)

var reportSearchFields = []string{"label", "text"}

// NewReport creates new instance of report data store.
func NewReport(c *Client) *ReportRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableReport); err != nil {
		log.Fatalf("could not create %s table: %s", tableReport, err)
	}

	if err := c.autoIndex(ctx, tableReport, dotagiftx2.Report{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableReport, err)
	}

	return &ReportRepository{c, reportSearchFields}
}

type ReportRepository struct {
	db            *Client
	keywordFields []string
}

func (s *ReportRepository) Find(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Report, error) {
	var res []dotagiftx2.Report
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *ReportRepository) Count(ctx context.Context, o dotagiftx2.FindOpts) (num int, err error) {
	o = dotagiftx2.FindOpts{
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
func (s *ReportRepository) includeRelatedFields(q r.Term) r.Term {
	return q.
		EqJoin(reportFieldUserID, r.Table(tableUser)).
		Map(func(t r.Term) r.Term {
			return t.Field("left").Merge(map[string]any{
				tableUser: t.Field("right"),
			})
		})
}

func (s *ReportRepository) Get(ctx context.Context, id string) (*dotagiftx2.Report, error) {
	row := &dotagiftx2.Report{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.ReportErrNotFound
		}

		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *ReportRepository) Create(ctx context.Context, in *dotagiftx2.Report) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *ReportRepository) table() r.Term {
	return r.Table(tableReport)
}
