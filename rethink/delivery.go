package rethink

import (
	"context"
	"errors"
	"log"

	"dario.cat/mergo"
	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableDelivery         = "delivery"
	deliveryFieldMarketID = "market_id"
)

var deliverySearchFields = []string{"id", "market_id"}

// NewDelivery creates new instance of delivery data store.
func NewDelivery(c *Client) *DeliveryRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableDelivery); err != nil {
		log.Fatalf("could not create %s table: %s", tableDelivery, err)
	}

	if err := c.autoIndex(ctx, tableDelivery, dotagiftx2.Delivery{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableDelivery, err)
	}

	return &DeliveryRepository{c, deliverySearchFields}
}

type DeliveryRepository struct {
	db            *Client
	keywordFields []string
}

func (s *DeliveryRepository) Find(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Delivery, error) {
	var res []dotagiftx2.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *DeliveryRepository) Count(ctx context.Context, o dotagiftx2.FindOpts) (num int, err error) {
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

func (s *DeliveryRepository) ToVerify(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Delivery, error) {
	var res []dotagiftx2.Delivery
	o.KeywordFields = s.keywordFields
	q := findOpts(o).parseOpts(s.table(), func(t r.Term) r.Term {
		return t.Filter(func(d r.Term) r.Term {
			return d.Field("retries").Default(0).Lt(dotagiftx2.DeliveryRetryLimit)
		})
	})
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

// includeRelatedFields injects user details base on market foreign keys.
func (s *DeliveryRepository) includeRelatedFields(q r.Term) r.Term {
	return q
	// return q.
	//	EqJoin(deliveryFieldMarketID, r.Table(tableMarket)).
	//	Map(func(t r.Term) r.Term {
	//		return t.Field("left").Merge(map[string]interface{}{
	//			tableMarket: t.Field("right"),
	//		})
	//	})
}

func (s *DeliveryRepository) Get(ctx context.Context, id string) (*dotagiftx2.Delivery, error) {
	row := &dotagiftx2.Delivery{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.DeliveryErrNotFound
		}

		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *DeliveryRepository) GetByMarketID(ctx context.Context, marketID string) (*dotagiftx2.Delivery, error) {
	var res []dotagiftx2.Delivery
	var err error

	q := s.table().GetAllByIndex(deliveryFieldMarketID, marketID)
	if err = s.db.list(ctx, q, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, dotagiftx2.DeliveryErrNotFound
	}

	return &res[0], nil
}

func (s *DeliveryRepository) Create(ctx context.Context, in *dotagiftx2.Delivery) error {
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

func (s *DeliveryRepository) Update(ctx context.Context, in *dotagiftx2.Delivery) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageMergeErr, err)
	}

	return nil
}

func (s *DeliveryRepository) table() r.Term {
	return r.Table(tableDelivery)
}
