package rethink

import (
	"context"
	"errors"
	"log"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableDelivery         = "delivery"
	deliveryFieldMarketID = "market_id"
)

var deliverySearchFields = []string{"id", "market_id"}

// NewDelivery creates new instance of delivery data store.
func NewDelivery(c *Client) *DeliveryStorage {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableDelivery); err != nil {
		log.Fatalf("could not create %s table: %s", tableDelivery, err)
	}

	if err := c.autoIndex(ctx, tableDelivery, dotagiftx.Delivery{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableDelivery, err)
	}

	return &DeliveryStorage{c, deliverySearchFields}
}

type DeliveryStorage struct {
	db            *Client
	keywordFields []string
}

func (s *DeliveryStorage) Find(ctx context.Context, o dotagiftx.FindOpts) ([]dotagiftx.Delivery, error) {
	var res []dotagiftx.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *DeliveryStorage) Count(ctx context.Context, o dotagiftx.FindOpts) (num int, err error) {
	o = dotagiftx.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
	}
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	err = s.db.one(ctx, q.Count(), &num)
	return
}

func (s *DeliveryStorage) ToVerify(ctx context.Context, o dotagiftx.FindOpts) ([]dotagiftx.Delivery, error) {
	var res []dotagiftx.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), func(t r.Term) r.Term {
		return t.Filter(func(d r.Term) r.Term {
			return d.Field("retries").Default(0).Lt(dotagiftx.DeliveryRetryLimit)
		})
	})
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

// includeRelatedFields injects user details base on market foreign keys.
func (s *DeliveryStorage) includeRelatedFields(q r.Term) r.Term {
	return q
	// return q.
	//	EqJoin(deliveryFieldMarketID, r.Table(tableMarket)).
	//	Map(func(t r.Term) r.Term {
	//		return t.Field("left").Merge(map[string]interface{}{
	//			tableMarket: t.Field("right"),
	//		})
	//	})
}

func (s *DeliveryStorage) Get(ctx context.Context, id string) (*dotagiftx.Delivery, error) {
	row := &dotagiftx.Delivery{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.DeliveryErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *DeliveryStorage) GetByMarketID(ctx context.Context, marketID string) (*dotagiftx.Delivery, error) {
	var res []dotagiftx.Delivery
	var err error

	q := s.table().GetAllByIndex(deliveryFieldMarketID, marketID)
	if err = s.db.list(ctx, q, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, dotagiftx.DeliveryErrNotFound
	}

	return &res[0], nil
}

func (s *DeliveryStorage) Create(ctx context.Context, in *dotagiftx.Delivery) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *DeliveryStorage) Update(ctx context.Context, in *dotagiftx.Delivery) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *DeliveryStorage) table() r.Term {
	return r.Table(tableDelivery)
}
