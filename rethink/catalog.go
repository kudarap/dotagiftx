package rethink

import (
	"log"
	"math"
	"time"

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
			qm := r.Table(tableMarket).
				Between(startTime, endTime, r.BetweenOpts{Index: marketFieldCreatedAt}).
				Filter(map[string]interface{}{marketFieldItemID: itemID})
			// Score rate evaluation.
			viewScore := t.Field(reductionField)
			entryScore := qm.Count()
			reserveScore := qm.Filter(map[string]interface{}{marketFieldStatus: core.MarketStatusReserved}).Count()
			soldScore := qm.Filter(map[string]interface{}{marketFieldStatus: core.MarketStatusSold}).Count()
			finalScore := r.Expr([]r.Term{
				viewScore.Mul(core.TrendScoreRateView),
				entryScore.Mul(core.TrendScoreRateMarketEntry),
				reserveScore.Mul(core.TrendScoreRateReserved),
				soldScore.Mul(core.TrendScoreRateSold),
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
	q := findOpts(o).parseOpts(s.table(), s.filterOutZeroQty)
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
	q := findOpts(o).parseOpts(s.table(), s.filterOutZeroQty)
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
	// Benchmark indexing.
	// avg proc time 2.555710726s
	tStart := time.Now()
	defer func() {
		s.logger.Infof("catalog indexed %s @ %s\n", itemID, time.Now().Sub(tStart))
	}()

	cat := &core.Catalog{}
	opts := core.FindOpts{Filter: core.Market{ItemID: itemID, Status: core.MarketStatusLive}}
	baseQ := newFindOptsQuery(r.Table(tableMarket), opts)

	var q r.Term
	var err error

	// Get item details by item ID.
	q = r.Table(tableItem).Get(itemID)
	if err = s.db.one(q, cat); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get total market count by item ID
	// and remove them if there's no entry
	quantity := baseQ.Count()
	if err = s.db.one(quantity, &cat.Quantity); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}
	if cat.Quantity == 0 {
		if err := s.zeroQtyCatalog(cat.ID); err != nil {
			return nil, errors.New(core.CatalogErrIndexing, err)
		}
	}

	// Get lowest sale price on the market by item ID.
	q = baseQ.Min("price").Field("price").Default(0)
	if err = s.db.one(q, &cat.LowestAsk); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get average sale price on the market by item ID.
	//q = baseQ.Avg("price").Default(0)
	//if err = s.db.one(q, &cat.AverageAsk); err != nil {
	//	return nil, errors.New(core.CatalogErrIndexing, err)
	//}

	skip := cat.Quantity / 2
	limit := 1
	if cat.Quantity%2 == 0 {
		limit = 2
		skip--
	} else {
		skip = int(math.Floor(float64(cat.Quantity) / 2))
	}
	//medianPos := quantity.Div(2).Floor()
	q = baseQ.OrderBy("price").Skip(skip).Limit(limit).Field("price").Default(0)
	if err = s.db.one(q, &cat.MeanAsk); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get highest price on the market by item ID.
	q = baseQ.Max("price").Field("price").Default(0)
	if err = s.db.one(q, &cat.HighestBid); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get recent_ask on the market by item ID.
	q = baseQ.Max("created_at").Field("created_at")
	recentAsk := &time.Time{}
	if err = s.db.one(q, recentAsk); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}
	cat.RecentAsk = recentAsk

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

func (s *catalogStorage) create(in *core.Catalog) error {
	// Fixes missing item in catalog that does not have views yet.
	in.ViewCount = 1
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	if _, err := s.db.insert(s.table().Insert(in)); err != nil {
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
	err = s.db.update(s.table().Get(in.ID).Update(in))
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
	q := s.table().Get(catalogID).Update(map[string]int{
		"quantity":    0,
		"lowest_ask":  0,
		"highest_bid": 0,
	})
	return s.db.update(q)
}

// NOTE! deprecated method and not being used.
func (s *catalogStorage) findIndex(o core.FindOpts) ([]core.Catalog, error) {
	q := s.indexBaseQuery()

	var res []core.Catalog
	o.KeywordFields = s.keywordFields
	q = newFindOptsQuery(q, o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *catalogStorage) indexBaseQuery() r.Term {
	return s.table().GroupByIndex(marketFieldItemID).Ungroup().
		Map(s.groupIndexMap).
		EqJoin(marketFieldItemID, r.Table(tableItem)).
		Zip()
}

func (s *catalogStorage) table() r.Term {
	return r.Table(tableCatalog)
}

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
