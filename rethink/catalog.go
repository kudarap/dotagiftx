package rethink

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/fatih/structs"
	"github.com/imdario/mergo"
	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableCatalog = "catalog"

// NewCatalog creates new instance of catalog data store.
func NewCatalog(c *Client, logger *logrus.Logger) core.CatalogStorage {
	if err := c.autoMigrate(tableCatalog); err != nil {
		log.Fatalf("could not create %s table: %s", tableCatalog, err)
	}

	if err := c.autoIndex(tableCatalog, core.Catalog{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableCatalog, err)
	}

	return &catalogStorage{c, itemSearchFields, logger}
}

type catalogStorage struct {
	db            *Client
	keywordFields []string
	logger        *logrus.Logger
}

func (s *catalogStorage) Trending() ([]core.Catalog, error) {
	/*
		r.db('dotagiftables')
		  .table('track')
		  .between(r.now().sub(604800), r.now(), {index: 'created_at'})
		  .filter({ type: 'v' })
		  .group('item_id').count()
		  .ungroup().orderBy(r.desc('reduction'))
		  .map(function(doc) {
		    let market = r.db('dotagiftables').table('market')
		        .between(r.now().sub(604800), r.now(), {index: 'created_at'})
		        .filter({item_id: doc('group')});

		    let viewScore = doc('reduction').mul(0.5);
		   	let entryScore = market.count().mul(0.1);
		    let reserveScore = market.filter({ status: 300 }).count().mul(4);
			let soldScore = market.filter({ status: 400 }).count().mul(4);
		    let score = r.expr([
		      viewScore,
		      entryScore,
		      reserveScore,
			  soldScore,
		    ]).sum();

		    return {
		      item_id: doc('group'),
		      score: score,
		      score_vw: viewScore,
		      score_ent: entryScore,
			  score_rsv: reserveScore,
		      score_sold: soldScore
		    }
		  })
		  .eqJoin('item_id', r.db('dotagiftables').table('catalog'))
		  .zip()
		  .orderBy(r.desc('score'))
		  .limit(10)
	*/

	// Scoring rate values from item views, entry, and reservations.
	const scoreFieldName = "score"

	// Date coverage for last 7 days.
	const last7Days = -time.Hour * 24 * 7
	endTime := time.Now()
	startTime := endTime.Add(last7Days)

	const reductionField = "reduction"
	q := r.Table(tableTrack).
		Between(startTime, endTime, r.BetweenOpts{Index: trackFieldCreatedAt}).
		Filter(map[string]string{trackFieldType: core.TrackTypeView}).
		Group(trackFieldItemID).Count().
		Ungroup().OrderBy(r.Desc(reductionField)).
		Map(func(t r.Term) interface{} {
			itemID := t.Field("group")
			mq := r.Table(tableMarket).
				Between(startTime, endTime, r.BetweenOpts{Index: marketFieldCreatedAt}).
				Filter(map[string]interface{}{marketFieldItemID: itemID})
			askQ := mq.Filter(map[string]interface{}{marketFieldType: core.MarketTypeAsk})
			// Score rate evaluation.
			viewScore := t.Field(reductionField)
			entryScore := askQ.Count()
			reserveScore := askQ.Filter(map[string]interface{}{marketFieldStatus: core.MarketStatusReserved}).Count()
			soldScore := askQ.Filter(map[string]interface{}{marketFieldStatus: core.MarketStatusSold}).Count()
			bidScore := mq.Filter(map[string]interface{}{marketFieldType: core.MarketTypeBid}).Count()
			finalScore := r.Expr([]r.Term{
				viewScore.Mul(core.TrendScoreRateView),
				entryScore.Mul(core.TrendScoreRateMarketEntry),
				reserveScore.Mul(core.TrendScoreRateReserved),
				soldScore.Mul(core.TrendScoreRateSold),
				bidScore.Mul(core.TrendScoreRateBid),
			}).Sum()

			return map[string]interface{}{
				trackFieldItemID: itemID,
				scoreFieldName:   finalScore,
				"view_count":     finalScore,
			}
		}).
		EqJoin(trackFieldItemID, r.Table(tableCatalog)).Zip().
		OrderBy(r.Desc(scoreFieldName)).Limit(10)

	var res []core.Catalog
	if err := s.db.list(q, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *catalogStorage) Find(o core.FindOpts) ([]core.Catalog, error) {
	var res []core.Catalog
	o.KeywordFields = s.keywordFields
	o.IndexSorting = true
	q := newFindOptsQuery(s.table(), o)
	//q := newCatalogFindOptsQuery(s.table(), o, s.filterOutZeroQty)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *catalogStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		IndexSorting:  true,
	}
	q := newFindOptsQuery(s.table(), o)
	//q := newCatalogFindOptsQuery(s.table(), o, s.filterOutZeroQty)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *catalogStorage) filterOutZeroQty(q r.Term) r.Term {
	return q.Filter(r.Row.Field("quantity").Gt(0))
}

func (s *catalogStorage) Get(id string) (*core.Catalog, error) {
	row, _ := s.getBySlug(id)
	if row != nil {
		return row, nil
	}

	row = &core.Catalog{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.CatalogErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *catalogStorage) getBySlug(slug string) (*core.Catalog, error) {
	row := &core.Catalog{}
	q := s.table().GetAllByIndex(itemFieldSlug, slug)
	if err := s.db.one(q, row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.CatalogErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *catalogStorage) Index(itemID string) (*core.Catalog, error) {
	bs := time.Now()
	defer func() {
		s.logger.Infof("catalog indexed %s @ %s\n", itemID, time.Now().Sub(bs))
	}()

	cat := &core.Catalog{}

	var q r.Term
	var err error

	// Get item details by item ID.
	q = r.Table(tableItem).Get(itemID)
	if err = s.db.one(q, cat); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get market offers from LIVE status.
	marketOffer := r.Table(tableMarket).Filter(core.Market{
		ItemID: itemID,
		Type:   core.MarketTypeAsk,
		Status: core.MarketStatusLive,
	})
	// Get offer count on the market by item ID.
	q = marketOffer.Count()
	if err = s.db.one(q, &cat.Quantity); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get ask count: %s", err))
	}
	if cat.Quantity > 0 {
		// Get lowest ask price on the market by item ID.
		q = marketOffer.Min(marketFieldPrice).Field(marketFieldPrice).Default(0)
		if err = s.db.one(q, &cat.LowestAsk); err != nil {
			return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get lowest ask price: %s", err))
		}

		// Get median ask price on the market by item ID.
		q = s.medianPriceQuery(cat.Quantity, marketOffer).Default(0)
		if err = s.db.one(q, &cat.MedianAsk); err != nil {
			return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get median ask price: %s", err))
		}

		// Get recent_ask on the market by item ID.
		q = marketOffer.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
		t := &time.Time{}
		if err = s.db.one(q, t); err != nil {
			return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get recent ask date: %s", err))
		}
		cat.RecentAsk = t
	}

	// Get market sales stats which calculated from RESERVED and SOLD statuses.
	marketSale := r.Table(tableMarket).Filter(core.Market{
		ItemID: itemID,
		Type:   core.MarketTypeAsk,
	}).Filter(func(doc r.Term) r.Term {
		return doc.Field(marketFieldStatus).Eq(core.MarketStatusReserved).
			Or(doc.Field(marketFieldStatus).Eq(core.MarketStatusSold))
	})
	// Get sale count on the market by item ID.
	q = marketSale.Count()
	if err = s.db.one(q, &cat.SaleCount); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get sales count: %s", err))
	}
	if cat.SaleCount > 0 {
		// Get average sale price on the market by item ID.
		q = marketSale.Avg(marketFieldPrice).Default(0)
		if err = s.db.one(q, &cat.AvgSale); err != nil {
			return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get avg sales price: %s", err))
		}
		// Get recent sale data on the market by item ID.
		q = marketSale.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
		t := &time.Time{}
		if err = s.db.one(q, t); err != nil {
			return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get recent sale date: %s", err))
		}
		cat.RecentSale = t
	}

	// Get reserved and sold count on the market by item ID.
	marketReservedCount := r.Table(tableMarket).Filter(core.Market{
		ItemID: itemID,
		Type:   core.MarketTypeAsk,
		Status: core.MarketStatusReserved,
	}).Count()
	if err = s.db.one(marketReservedCount, &cat.ReservedCount); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, fmt.Errorf("could not get reserved count: %s", err))
	}
	cat.SoldCount = cat.SaleCount - cat.ReservedCount

	// Check for exiting entry for update or create.
	if cur, _ := s.Get(itemID); cur == nil {
		err = s.create(cat)
	} else {
		err = s.update(cat)
	}

	if err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	return cat, nil
}

func (s *catalogStorage) medianPriceQuery(qty int, t r.Term) r.Term {
	q := t.OrderBy(marketFieldPrice)
	if qty < 2 {
		return q.Field(marketFieldPrice)
	}

	skip := int(math.Floor(float64(qty) / 2))
	limit := 1
	if qty%2 == 0 {
		skip--
		limit = 2
	}

	return q.Skip(skip).Limit(limit).Avg(marketFieldPrice)
}

func (s *catalogStorage) create(in *core.Catalog) error {
	// Fixes missing item in catalog that does not have views yet.
	in.ViewCount = 1
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	// Convert catalog into map to insert zero value fields.
	m := catalogToMap(in)

	if _, err := s.db.insert(s.table().Insert(m)); err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	return nil
}

func (s *catalogStorage) update(in *core.Catalog) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	// Convert catalog into map to insert zero value fields.
	m := catalogToMap(in)

	err = s.db.update(s.table().Get(in.ID).Update(m))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return errors.New(core.StorageMergeErr, err)
	}

	return nil
}

// zeroQtyCatalog reset the catalog entry price when it reaches zero entry/qty.
func (s *catalogStorage) zeroQtyCatalog(catalogID string) error {
	cat := map[string]interface{}{
		"quantity":   0,
		"lowest_ask": 0,
		"median_ask": 0,
		//"highest_bid": 0,
		"recent_ask": nil,
	}

	var err error
	if cur, _ := s.Get(catalogID); cur == nil {
		_, err = s.db.insert(s.table().Insert(cat))
	} else {
		err = s.db.update(s.table().Get(catalogID).Update(cat))
	}

	return err
}

func (s *catalogStorage) table() r.Term {
	return r.Table(tableCatalog)
}

func catalogToMap(cat *core.Catalog) map[string]interface{} {
	s := structs.New(cat)
	s.TagName = "json"
	return s.Map()
}

// NOTE! deprecated method and not being used and for reference only.
func (s *catalogStorage) findIndexLegacy(o core.FindOpts) ([]core.Catalog, error) {
	q := s.indexBaseQuery()

	var res []core.Catalog
	o.KeywordFields = s.keywordFields
	q = newFindOptsQuery(q, o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

// NOTE! deprecated method and not being used and for reference only.
func (s *catalogStorage) indexBaseQuery() r.Term {
	return s.table().GroupByIndex(marketFieldItemID).Ungroup().
		Map(s.groupIndexMap).
		EqJoin(marketFieldItemID, r.Table(tableItem)).
		Zip()
}

// NOTE! deprecated method and not being used and for reference only.
func (s *catalogStorage) groupIndexMap(catalog r.Term) interface{} {
	//r.db('dotagiftables').table('market').group({index: 'item_id'}).ungroup().map(
	//    function (doc) {
	//      let liveMarket = doc('reduction').filter({status: 200});
	//      return {
	//        item_id: doc('group'),
	//        quantity: liveMarket.count(),
	//        lowest_ask: liveMarket.min('price')('price').default(0),
	//        highest_bid: liveMarket.max('price')('price').default(0),
	//        recent_ask: liveMarket.max('created_at')('created_at').default(null),
	//        item: r.db('dotagiftables').table('item').get(doc('group')),
	//      };
	//    }
	//)

	id := catalog.Field("group")
	live := catalog.Field("reduction").Filter(core.Market{Status: core.MarketStatusLive})
	return struct {
		ItemID     r.Term `db:"item_id"`
		Quantity   r.Term `db:"quantity"`
		LowestAsk  r.Term `db:"lowest_ask"`
		HighestBid r.Term `db:"highest_bid"`
		RecentAsk  r.Term `db:"recent_ask"`
		//Item       r.Term `db:"item"`
	}{
		id,
		live.Count().Default(0),
		live.Min("price").Field("price").Default(0),
		live.Max("price").Field("price").Default(0),
		live.Max("created_at").Field("created_at").Default(nil),
		//r.Table(tableItem).Get(id),
	}
}

// NOTE! deprecated method and not being used and for reference only.
func (s *catalogStorage) trendingV0() ([]core.Catalog, error) {
	/*
		r.db('d2g')
		.table('track')
		.filter({type: 'v'})
		.orderBy(r.desc('created_at'))
		.limit(100)
		.group('item_id').count()
		.ungroup().orderBy(r.desc('reduction'))
		.map(function(doc) {
		  return {
		    item_id: doc('group'),
		    score: doc('reduction'),
		  }
		})
		.eqJoin('item_id', r.db('d2g').table('catalog'))
		.zip()
		.orderBy(r.desc('score'))
		.limit(10)
	*/

	// Accumulate views of the recent 100 records.
	q := r.Table(tableTrack).Filter(map[string]string{"type": "v"}).
		OrderBy(r.Desc("created_at")).Limit(100).
		Group("item_id").Count().
		Ungroup().OrderBy(r.Desc("reduction")).
		Map(func(t r.Term) interface{} {
			return map[string]interface{}{
				"item_id": t.Field("group"),
				"score":   t.Field("reduction"),
			}
		}).
		EqJoin("item_id", r.Table(tableCatalog)).Zip().
		OrderBy(r.Desc("score")).Limit(10)

	var res []core.Catalog
	if err := s.db.list(q, &res); err != nil {
		return nil, err
	}

	return res, nil
}
