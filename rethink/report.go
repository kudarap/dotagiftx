package rethink

import (
	"log"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableReport = "report"

var reportSearchFields = []string{"label", "text"}

// NewReport creates new instance of report data store.
func NewReport(c *Client) core.ReportStorage {
	if err := c.autoMigrate(tableReport); err != nil {
		log.Fatalf("could not create %s table: %s", tableReport, err)
	}

	if err := c.autoIndex(tableCatalog, core.Report{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableCatalog, err)
	}

	return &reportStorage{c, reportSearchFields}
}

type reportStorage struct {
	db            *Client
	keywordFields []string
}

func (s *reportStorage) Find(o core.FindOpts) ([]core.Report, error) {
	var res []core.Report
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *reportStorage) Count(o core.FindOpts) (num int, err error) {
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

func (s *reportStorage) Get(id string) (*core.Report, error) {
	row := &core.Report{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.ReportErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *reportStorage) Create(in *core.Report) error {
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

func (s *reportStorage) table() r.Term {
	return r.Table(tableReport)
}
