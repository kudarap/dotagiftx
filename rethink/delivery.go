package rethink

import (
	"log"

	"dario.cat/mergo"
	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableDelivery         = "delivery"
	deliveryFieldMarketID = "market_id"
)

var deliverySearchFields = []string{"id", "market_id"}

// NewDelivery creates new instance of delivery data store.
func NewDelivery(c *Client) dgx.DeliveryStorage {
	if err := c.autoMigrate(tableDelivery); err != nil {
		log.Fatalf("could not create %s table: %s", tableDelivery, err)
	}

	if err := c.autoIndex(tableDelivery, dgx.Delivery{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableDelivery, err)
	}

	return &deliveryStorage{c, deliverySearchFields}
}

type deliveryStorage struct {
	db            *Client
	keywordFields []string
}

func (s *deliveryStorage) Find(o dgx.FindOpts) ([]dgx.Delivery, error) {
	var res []dgx.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *deliveryStorage) Count(o dgx.FindOpts) (num int, err error) {
	o = dgx.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
	}
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *deliveryStorage) ToVerify(o dgx.FindOpts) ([]dgx.Delivery, error) {
	var res []dgx.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), func(t r.Term) r.Term {
		return t.Filter(func(d r.Term) r.Term {
			return d.Field("retries").Default(0).Lt(dgx.DeliveryRetryLimit)
		})
	})
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return res, nil
}

// includeRelatedFields injects user details base on market foreign keys.
func (s *deliveryStorage) includeRelatedFields(q r.Term) r.Term {
	return q
	//return q.
	//	EqJoin(deliveryFieldMarketID, r.Table(tableMarket)).
	//	Map(func(t r.Term) r.Term {
	//		return t.Field("left").Merge(map[string]interface{}{
	//			tableMarket: t.Field("right"),
	//		})
	//	})
}

func (s *deliveryStorage) Get(id string) (*dgx.Delivery, error) {
	row := &dgx.Delivery{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.DeliveryErrNotFound
		}

		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *deliveryStorage) GetByMarketID(marketID string) (*dgx.Delivery, error) {
	var res []dgx.Delivery
	var err error

	q := s.table().GetAllByIndex(deliveryFieldMarketID, marketID)
	if err = s.db.list(q, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, dgx.DeliveryErrNotFound
	}

	return &res[0], nil
}

func (s *deliveryStorage) Create(in *dgx.Delivery) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *deliveryStorage) Update(in *dgx.Delivery) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return errors.New(dgx.StorageMergeErr, err)
	}

	return nil
}

func (s *deliveryStorage) table() r.Term {
	return r.Table(tableDelivery)
}
