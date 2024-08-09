package rethink

import (
	"fmt"
	"math"
	"time"

	"dario.cat/mergo"
	"github.com/fatih/structs"
	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/log"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableCatalog = "catalog"

// NewCatalog creates new instance of catalog data store.
func NewCatalog(c *Client, lg log.Logger) dgx.CatalogStorage {
	if err := c.autoMigrate(tableCatalog); err != nil {
		lg.Fatalf("could not create %s table: %s", tableCatalog, err)
	}

	if err := c.autoIndex(tableCatalog, dgx.Catalog{}); err != nil {
		lg.Fatalf("could not create index on %s table: %s", tableCatalog, err)
	}

	return &catalogStorage{c, itemSearchFields, lg}
}

type catalogStorage struct {
	db            *Client
	keywordFields []string
	logger        log.Logger
}

func (s *catalogStorage) Trending() ([]dgx.Catalog, error) {
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
		    let bidScore = market.filter({ type: 20 }).count().mul(2);
		    let score = r.expr([
		      viewScore,
		      entryScore,
		      reserveScore,
		      soldScore,
		      bidScore,
		    ]).sum();

		    return {
		      item_id: doc('group'),
		      score: score,
		      score_vw: viewScore,
		      score_ent: entryScore,
		      score_rsv: reserveScore,
		      score_sold: soldScore,
		      score_bid: bidScore
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
		Filter(map[string]string{trackFieldType: dgx.TrackTypeView}).
		Group(trackFieldItemID).Count().
		Ungroup().OrderBy(r.Desc(reductionField)).
		Map(func(t r.Term) interface{} {
			itemID := t.Field("group")
			mq := r.Table(tableMarket).
				Between(startTime, endTime, r.BetweenOpts{Index: marketFieldCreatedAt}).
				Filter(map[string]interface{}{marketFieldItemID: itemID})
			askQ := mq.Filter(map[string]interface{}{marketFieldType: dgx.MarketTypeAsk})
			// Score rate evaluation.
			viewScore := t.Field(reductionField)
			entryScore := askQ.Count()
			reserveScore := askQ.Filter(map[string]interface{}{marketFieldStatus: dgx.MarketStatusReserved}).Count()
			soldScore := askQ.Filter(map[string]interface{}{marketFieldStatus: dgx.MarketStatusSold}).Count()
			bidScore := mq.Filter(map[string]interface{}{marketFieldType: dgx.MarketTypeBid}).Count()
			finalScore := r.Expr([]r.Term{
				viewScore.Mul(dgx.TrendScoreRateView),
				entryScore.Mul(dgx.TrendScoreRateMarketEntry),
				reserveScore.Mul(dgx.TrendScoreRateReserved),
				soldScore.Mul(dgx.TrendScoreRateSold),
				bidScore.Mul(dgx.TrendScoreRateBid),
			}).Sum()

			return map[string]interface{}{
				trackFieldItemID: itemID,
				scoreFieldName:   finalScore,
				"view_count":     finalScore,
			}
		}).
		EqJoin(trackFieldItemID, r.Table(tableCatalog)).Zip().
		OrderBy(r.Desc(scoreFieldName)).Limit(10)

	var res []dgx.Catalog
	if err := s.db.list(q, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *catalogStorage) Find(o dgx.FindOpts) ([]dgx.Catalog, error) {
	var res []dgx.Catalog
	o.KeywordFields = s.keywordFields
	o.IndexSorting = true
	q := newFindOptsQuery(s.table(), o)
	//q := newCatalogFindOptsQuery(s.table(), o, s.filterOutZeroQty)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}
	return res, nil
}

func (s *catalogStorage) Count(o dgx.FindOpts) (num int, err error) {
	o = dgx.FindOpts{
		KeywordFields: s.keywordFields,
		IndexSorting:  true,
		Keyword:       o.Keyword,
		Filter:        o.Filter,
		Sort:          o.Sort,
	}
	q := newFindOptsQuery(s.table(), o)
	//q := newCatalogFindOptsQuery(s.table(), o, s.filterOutZeroQty)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *catalogStorage) filterOutZeroQty(q r.Term) r.Term {
	return q.Filter(r.Row.Field("quantity").Gt(0))
}

func (s *catalogStorage) Get(id string) (*dgx.Catalog, error) {
	row, _ := s.getBySlug(id)
	if row != nil {
		return row, nil
	}

	row = &dgx.Catalog{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.CatalogErrNotFound
		}

		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *catalogStorage) getBySlug(slug string) (*dgx.Catalog, error) {
	row := &dgx.Catalog{}
	q := s.table().GetAllByIndex(itemFieldSlug, slug)
	if err := s.db.one(q, row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.CatalogErrNotFound
		}

		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *catalogStorage) Index(itemID string) (*dgx.Catalog, error) {
	bs := time.Now()
	defer func() {
		s.logger.Infof("catalog indexed %s @ %s\n", itemID, time.Now().Sub(bs))
	}()

	var benchStart time.Time

	cat := &dgx.Catalog{}

	var q r.Term
	var err error

	// Get item details by item ID.
	q = r.Table(tableItem).Get(itemID)
	if err = s.db.one(q, cat); err != nil {
		return nil, errors.New(dgx.CatalogErrIndexing, err)
	}

	benchStart = time.Now()
	// Get market offers summary from LIVE status.
	cat.Quantity, cat.LowestAsk, cat.MedianAsk, cat.RecentAsk, err = s.getOffersSummary(itemID)
	if err != nil {
		return nil, errors.New(dgx.CatalogErrIndexing, err)
	}
	s.logger.Println("rethink/catalog getOffersSummary", time.Now().Sub(benchStart))

	benchStart = time.Now()
	// Get market buy orders summary.
	cat.BidCount, cat.HighestBid, cat.RecentBid, err = s.getBuyOrdersSummary(itemID)
	if err != nil {
		return nil, errors.New(dgx.CatalogErrIndexing, err)
	}
	s.logger.Println("rethink/catalog getBuyOrdersSummary", time.Now().Sub(benchStart))

	benchStart = time.Now()
	// Get market sales stats which calculated from RESERVED and SOLD statuses.
	cat.SaleCount, cat.AvgSale, cat.RecentSale, err = s.getSaleSummary(itemID)
	if err != nil {
		return nil, errors.New(dgx.CatalogErrIndexing, err)
	}
	s.logger.Println("rethink/catalog getSaleSummary", time.Now().Sub(benchStart))

	benchStart = time.Now()
	// Get reserved and sold count on the market by item ID.
	cat.ReservedCount, err = s.getReservedCounts(itemID)
	if err != nil {
		return nil, errors.New(dgx.CatalogErrIndexing, err)
	}
	cat.SoldCount = cat.SaleCount - cat.ReservedCount
	s.logger.Println("rethink/catalog getReservedCounts", time.Now().Sub(benchStart))

	// Check for exiting entry for update or create.
	if cur, _ := s.Get(itemID); cur == nil {
		err = s.create(cat)
	} else {
		err = s.update(cat)
	}

	if err != nil {
		return nil, errors.New(dgx.CatalogErrIndexing, err)
	}

	return cat, nil
}

// getOffersSummary returns market offers summary from LIVE status.
func (s *catalogStorage) getOffersSummary(itemID string) (count int, lowest, median float64, recent *time.Time, err error) {
	// Get market offers from LIVE status.
	offer := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dgx.Market{
			Type:            dgx.MarketTypeAsk,
			Status:          dgx.MarketStatusLive,
			InventoryStatus: dgx.InventoryStatusVerified,
		})
	// Get offer count on the market by item ID.
	q := offer.Count()
	if err = s.db.one(q, &count); err != nil {
		err = fmt.Errorf("could not get ask count: %s", err)
		return
	}
	if count == 0 {
		return
	}

	// Get the lowest ask price on the market by item ID.
	q = offer.Min(marketFieldPrice).Field(marketFieldPrice).Default(0)
	if err = s.db.one(q, &lowest); err != nil {
		err = fmt.Errorf("could not get lowest ask price: %s", err)
		return
	}

	// Get median ask price on the market by item ID.
	q = s.medianPriceQuery(count, offer).Default(0)
	if err = s.db.one(q, &median); err != nil {
		err = fmt.Errorf("could not get median ask price: %s", err)
		return
	}

	// Get recent ask on the market by item ID.
	q = offer.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
	t := &time.Time{}
	if err = s.db.one(q, t); err != nil {
		err = fmt.Errorf("could not get recent ask date: %s", err)
		return
	}
	recent = t
	return
}

// getBuyOrdersSummary returns market buy orders from BID type and LIVE status.
func (s *catalogStorage) getBuyOrdersSummary(itemID string) (count int, max float64, recent *time.Time, err error) {
	buyOrder := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dgx.Market{
			Type:   dgx.MarketTypeBid,
			Status: dgx.MarketStatusLive,
		})
	// Get bid count on the market by item ID.
	q := buyOrder.Count()
	if err = s.db.one(q, &count); err != nil {
		err = fmt.Errorf("could not get bid count: %s", err)
		return
	}
	if count == 0 {
		return
	}

	// Get the highest bid price on the market by item ID.
	q = buyOrder.Max(marketFieldPrice).Field(marketFieldPrice).Default(0)
	if err = s.db.one(q, &max); err != nil {
		err = fmt.Errorf("could not get highest bid price: %s", err)
		return
	}

	// Get recent bid on the market by item ID.
	q = buyOrder.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
	t := &time.Time{}
	if err = s.db.one(q, t); err != nil {
		err = fmt.Errorf("could not get recent bid date: %s", err)
		return
	}
	recent = t
	return
}

// getSaleSummary returns market sales stats which calculated from RESERVED and SOLD statuses.
func (s *catalogStorage) getSaleSummary(itemID string) (count int, avg float64, recent *time.Time, err error) {
	sale := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dgx.Market{
			Type: dgx.MarketTypeAsk,
		}).Filter(func(doc r.Term) r.Term {
		return doc.Field(marketFieldStatus).Eq(dgx.MarketStatusReserved).
			Or(doc.Field(marketFieldStatus).Eq(dgx.MarketStatusSold))
	})

	// Get sale count on the market by item ID.
	q := sale.Count()
	if err = s.db.one(q, &count); err != nil {
		err = fmt.Errorf("could not get sales count: %s", err)
		return
	}
	if count == 0 {
		return
	}

	// Get average sale price on the market by item ID.
	q = sale.Avg(marketFieldPrice).Default(0)
	if err = s.db.one(q, &avg); err != nil {
		err = fmt.Errorf("could not get avg sales price: %s", err)
		return
	}
	// Get recent sale data on the market by item ID.
	q = sale.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
	t := &time.Time{}
	if err = s.db.one(q, t); err != nil {
		err = fmt.Errorf("could not get recent sale date: %s", err)
		return
	}
	recent = t
	return
}

func (s *catalogStorage) getReservedCounts(itemID string) (count int, err error) {
	reserved := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dgx.Market{
			Type:   dgx.MarketTypeAsk,
			Status: dgx.MarketStatusReserved,
		}).Count()

	if err = s.db.one(reserved, &count); err != nil {
		err = fmt.Errorf("could not get reserved count: %s", err)
	}

	return
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

func (s *catalogStorage) create(in *dgx.Catalog) error {
	// Fixes missing item in catalog that does not have views yet.
	in.ViewCount = 1
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	// Convert catalog into map to insert zero value fields.
	m := catalogToMap(in)

	if _, err := s.db.insert(s.table().Insert(m)); err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}

	return nil
}

func (s *catalogStorage) update(in *dgx.Catalog) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	// Convert catalog into map to insert zero value fields.
	m := catalogToMap(in)

	err = s.db.update(s.table().Get(in.ID).Update(m))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return errors.New(dgx.StorageMergeErr, err)
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

func catalogToMap(cat *dgx.Catalog) map[string]interface{} {
	s := structs.New(cat)
	s.TagName = "json"
	return s.Map()
}

// NOTE! deprecated method and not being used and for reference only.
func (s *catalogStorage) findIndexLegacy(o dgx.FindOpts) ([]dgx.Catalog, error) {
	q := s.indexBaseQuery()

	var res []dgx.Catalog
	o.KeywordFields = s.keywordFields
	q = newFindOptsQuery(q, o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
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
	live := catalog.Field("reduction").Filter(dgx.Market{Status: dgx.MarketStatusLive})
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
func (s *catalogStorage) trendingV0() ([]dgx.Catalog, error) {
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

	var res []dgx.Catalog
	if err := s.db.list(q, &res); err != nil {
		return nil, err
	}

	return res, nil
}
