package rethink

import (
	"log"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableReport       = "report"
	reportFieldUserID = "user_id"
)

var reportSearchFields = []string{"label", "text"}

// NewReport creates new instance of report data store.
func NewReport(c *Client) dotagiftx.ReportStorage {
	if err := c.autoMigrate(tableReport); err != nil {
		log.Fatalf("could not create %s table: %s", tableReport, err)
	}

	if err := c.autoIndex(tableReport, dotagiftx.Report{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableReport, err)
	}

	return &reportStorage{c, reportSearchFields}
}

type reportStorage struct {
	db            *Client
	keywordFields []string
}

func (s *reportStorage) Find(o dotagiftx.FindOpts) ([]dotagiftx.Report, error) {
	var res []dotagiftx.Report
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *reportStorage) Count(o dotagiftx.FindOpts) (num int, err error) {
	o = dotagiftx.FindOpts{
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
func (s *reportStorage) includeRelatedFields(q r.Term) r.Term {
	return q.
		EqJoin(reportFieldUserID, r.Table(tableUser)).
		Map(func(t r.Term) r.Term {
			return t.Field("left").Merge(map[string]interface{}{
				tableUser: t.Field("right"),
			})
		})
}

func (s *reportStorage) Get(id string) (*dotagiftx.Report, error) {
	row := &dotagiftx.Report{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dotagiftx.ReportErrNotFound
		}

		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *reportStorage) Create(in *dotagiftx.Report) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *reportStorage) table() r.Term {
	return r.Table(tableReport)
}
