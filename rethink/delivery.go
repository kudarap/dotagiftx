package rethink

import (
	"log"

	"github.com/imdario/mergo"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableDelivery       = "delivery"
	deliveryFieldUserID = "user_id"
)

var deliverySearchFields = []string{"label", "text"}

// NewDelivery creates new instance of delivery data store.
func NewDelivery(c *Client) core.DeliveryStorage {
	if err := c.autoMigrate(tableDelivery); err != nil {
		log.Fatalf("could not create %s table: %s", tableDelivery, err)
	}

	if err := c.autoIndex(tableDelivery, core.Delivery{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableDelivery, err)
	}

	return &deliveryStorage{c, deliverySearchFields}
}

type deliveryStorage struct {
	db            *Client
	keywordFields []string
}

func (s *deliveryStorage) Find(o core.FindOpts) ([]core.Delivery, error) {
	var res []core.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *deliveryStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{
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
func (s *deliveryStorage) includeRelatedFields(q r.Term) r.Term {
	return q.
		EqJoin(deliveryFieldUserID, r.Table(tableUser)).
		Map(func(t r.Term) r.Term {
			return t.Field("left").Merge(map[string]interface{}{
				tableUser: t.Field("right"),
			})
		})
}

func (s *deliveryStorage) Get(id string) (*core.Delivery, error) {
	row := &core.Delivery{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.DeliveryErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *deliveryStorage) Create(in *core.Delivery) error {
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

func (s *deliveryStorage) Update(in *core.Delivery) error {
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

func (s *deliveryStorage) table() r.Term {
	return r.Table(tableDelivery)
}
