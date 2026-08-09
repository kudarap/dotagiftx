package rethink

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"
	"time"

	"dario.cat/mergo"
	"github.com/fatih/structs"
	"github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableCatalog = "catalog"

// NewCatalog creates new instance of catalog data store.
func NewCatalog(c *Client, lg *slog.Logger) *CatalogRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableCatalog); err != nil {
		lg.ErrorContext(ctx, "could not create catalog table", "table", tableCatalog, "error", err)
		os.Exit(1)
	}

	if err := c.autoIndex(ctx, tableCatalog, dotagiftx.Catalog{}); err != nil {
		lg.ErrorContext(ctx, "could not create index on catalog table", "table", tableCatalog, "error", err)
		os.Exit(1)
	}

	return &CatalogRepository{c, itemSearchFields, lg}
}

type CatalogRepository struct {
	db            *Client
	keywordFields []string
	logger        *slog.Logger
}

func (s *CatalogRepository) Trending(ctx context.Context) ([]dotagiftx.Catalog, error) {
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
		Filter(map[string]string{trackFieldType: dotagiftx.TrackTypeView}).
		Group(trackFieldItemID).Count().
		Ungroup().OrderBy(r.Desc(reductionField)).
		Map(func(t r.Term) any {
			itemID := t.Field("group")
			mq := r.Table(tableMarket).
				Between(startTime, endTime, r.BetweenOpts{Index: marketFieldCreatedAt}).
				Filter(map[string]any{marketFieldItemID: itemID})
			askQ := mq.Filter(map[string]any{marketFieldType: dotagiftx.MarketTypeAsk})
			// Score rate evaluation.
			viewScore := t.Field(reductionField)
			entryScore := askQ.Count()
			reserveScore := askQ.Filter(map[string]any{marketFieldStatus: dotagiftx.MarketStatusReserved}).Count()
			soldScore := askQ.Filter(map[string]any{marketFieldStatus: dotagiftx.MarketStatusSold}).Count()
			bidScore := mq.Filter(map[string]any{marketFieldType: dotagiftx.MarketTypeBid}).Count()
			finalScore := r.Expr([]r.Term{
				viewScore.Mul(dotagiftx.TrendScoreRateView),
				entryScore.Mul(dotagiftx.TrendScoreRateMarketEntry),
				reserveScore.Mul(dotagiftx.TrendScoreRateReserved),
				soldScore.Mul(dotagiftx.TrendScoreRateSold),
				bidScore.Mul(dotagiftx.TrendScoreRateBid),
			}).Sum()

			return map[string]any{
				trackFieldItemID: itemID,
				scoreFieldName:   finalScore,
				"view_count":     finalScore,
			}
		}).
		EqJoin(trackFieldItemID, r.Table(tableCatalog)).Zip().
		OrderBy(r.Desc(scoreFieldName)).Limit(10)

	var res []dotagiftx.Catalog
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CatalogRepository) Find(ctx context.Context, o dotagiftx.FindOpts) ([]dotagiftx.Catalog, error) {
	var res []dotagiftx.Catalog
	o.KeywordFields = s.keywordFields
	o.IndexSorting = true
	q := newFindOptsQuery(s.table(), o)
	// q := newCatalogFindOptsQuery(s.table(), o, s.filterOutZeroQty)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}
	return res, nil
}

func (s *CatalogRepository) Count(ctx context.Context, o dotagiftx.FindOpts) (num int, err error) {
	o = dotagiftx.FindOpts{
		KeywordFields: s.keywordFields,
		IndexSorting:  true,
		Keyword:       o.Keyword,
		Filter:        o.Filter,
		Sort:          o.Sort,
	}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(ctx, q.Count(), &num)
	return
}

func (s *CatalogRepository) Get(ctx context.Context, id string) (*dotagiftx.Catalog, error) {
	row, _ := s.getBySlug(ctx, id)
	if row != nil {
		return row, nil
	}

	row = &dotagiftx.Catalog{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.CatalogErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *CatalogRepository) getBySlug(ctx context.Context, slug string) (*dotagiftx.Catalog, error) {
	row := &dotagiftx.Catalog{}
	q := s.table().GetAllByIndex(itemFieldSlug, slug)
	if err := s.db.one(ctx, q, row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.CatalogErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *CatalogRepository) Index(ctx context.Context, itemID string) (*dotagiftx.Catalog, error) {
	bs := time.Now()
	defer func() {
		s.logger.InfoContext(ctx, "catalog indexed", "item_id", itemID, "elapsed", time.Since(bs))
	}()

	var benchStart time.Time

	cat := &dotagiftx.Catalog{}

	var q r.Term
	var err error

	// Get item details by item ID.
	q = r.Table(tableItem).Get(itemID)
	if err = s.db.one(ctx, q, cat); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.CatalogErrIndexing, err)
	}

	benchStart = time.Now()
	// Get market offers summary from LIVE status.
	cat.Quantity, cat.LowestAsk, cat.MedianAsk, cat.RecentAsk, err = s.getOffersSummary(ctx, itemID)
	if err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.CatalogErrIndexing, err)
	}
	s.logger.InfoContext(ctx, "rethink/catalog getOffersSummary", "elapsed", time.Since(benchStart))

	benchStart = time.Now()
	// Get market buy orders summary.
	cat.BidCount, cat.HighestBid, cat.RecentBid, err = s.getBuyOrdersSummary(ctx, itemID)
	if err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.CatalogErrIndexing, err)
	}
	s.logger.InfoContext(ctx, "rethink/catalog getBuyOrdersSummary", "elapsed", time.Since(benchStart))

	benchStart = time.Now()
	// Get market sales stats which calculated from RESERVED and SOLD statuses.
	cat.SaleCount, cat.AvgSale, cat.RecentSale, err = s.getSaleSummary(ctx, itemID)
	if err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.CatalogErrIndexing, err)
	}
	s.logger.InfoContext(ctx, "rethink/catalog getSaleSummary", "elapsed", time.Since(benchStart))

	benchStart = time.Now()
	// Get reserved and sold count on the market by item ID.
	cat.ReservedCount, err = s.getReservedCounts(ctx, itemID)
	if err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.CatalogErrIndexing, err)
	}
	cat.SoldCount = cat.SaleCount - cat.ReservedCount
	s.logger.InfoContext(ctx, "rethink/catalog getReservedCounts", "elapsed", time.Since(benchStart))

	// Check for exiting entry for update or create.
	if cur, _ := s.Get(ctx, itemID); cur == nil {
		err = s.create(ctx, cat)
	} else {
		err = s.update(ctx, cat)
	}

	if err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.CatalogErrIndexing, err)
	}

	return cat, nil
}

// getOffersSummary returns market offers summary from LIVE status.
func (s *CatalogRepository) getOffersSummary(ctx context.Context, itemID string) (count int, lowest, median float64, recent *time.Time, err error) {
	// Get market offers from LIVE status.
	offer := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dotagiftx.Market{
			Type:            dotagiftx.MarketTypeAsk,
			Status:          dotagiftx.MarketStatusLive,
			InventoryStatus: dotagiftx.InventoryStatusVerified,
		})
	// Get offer count on the market by item ID.
	q := offer.Count()
	if err = s.db.one(ctx, q, &count); err != nil {
		err = fmt.Errorf("could not get ask count: %s", err)
		return
	}
	if count == 0 {
		return
	}

	// Get the lowest ask price on the market by item ID.
	q = offer.Min(marketFieldPrice).Field(marketFieldPrice).Default(0)
	if err = s.db.one(ctx, q, &lowest); err != nil {
		err = fmt.Errorf("could not get lowest ask price: %s", err)
		return
	}

	// Get median ask price on the market by item ID.
	q = s.medianPriceQuery(count, offer).Default(0)
	if err = s.db.one(ctx, q, &median); err != nil {
		err = fmt.Errorf("could not get median ask price: %s", err)
		return
	}

	// Get recent ask on the market by item ID.
	q = offer.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
	t := &time.Time{}
	if err = s.db.one(ctx, q, t); err != nil {
		err = fmt.Errorf("could not get recent ask date: %s", err)
		return
	}
	recent = t
	return
}

// getBuyOrdersSummary returns market buy orders from BID type and LIVE status.
func (s *CatalogRepository) getBuyOrdersSummary(ctx context.Context, itemID string) (count int, max float64, recent *time.Time, err error) {
	buyOrder := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dotagiftx.Market{
			Type:   dotagiftx.MarketTypeBid,
			Status: dotagiftx.MarketStatusLive,
		})
	// Get bid count on the market by item ID.
	q := buyOrder.Count()
	if err = s.db.one(ctx, q, &count); err != nil {
		err = fmt.Errorf("could not get bid count: %s", err)
		return
	}
	if count == 0 {
		return
	}

	// Get the highest bid price on the market by item ID.
	q = buyOrder.Max(marketFieldPrice).Field(marketFieldPrice).Default(0)
	if err = s.db.one(ctx, q, &max); err != nil {
		err = fmt.Errorf("could not get highest bid price: %s", err)
		return
	}

	// Get recent bid on the market by item ID.
	q = buyOrder.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
	t := &time.Time{}
	if err = s.db.one(ctx, q, t); err != nil {
		err = fmt.Errorf("could not get recent bid date: %s", err)
		return
	}
	recent = t
	return
}

// getSaleSummary returns market sales stats which calculated from RESERVED and SOLD statuses.
func (s *CatalogRepository) getSaleSummary(ctx context.Context, itemID string) (count int, avg float64, recent *time.Time, err error) {
	sale := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dotagiftx.Market{
			Type: dotagiftx.MarketTypeAsk,
		}).Filter(func(doc r.Term) r.Term {
		return doc.Field(marketFieldStatus).Eq(dotagiftx.MarketStatusReserved).
			Or(doc.Field(marketFieldStatus).Eq(dotagiftx.MarketStatusSold))
	})

	// Get sale count on the market by item ID.
	q := sale.Count()
	if err = s.db.one(ctx, q, &count); err != nil {
		err = fmt.Errorf("could not get sales count: %s", err)
		return
	}
	if count == 0 {
		return
	}

	// Get average sale price on the market by item ID.
	q = sale.Avg(marketFieldPrice).Default(0)
	if err = s.db.one(ctx, q, &avg); err != nil {
		err = fmt.Errorf("could not get avg sales price: %s", err)
		return
	}
	// Get recent sale data on the market by item ID.
	q = sale.Max(marketFieldCreatedAt).Field(marketFieldCreatedAt).Default(nil)
	t := &time.Time{}
	if err = s.db.one(ctx, q, t); err != nil {
		err = fmt.Errorf("could not get recent sale date: %s", err)
		return
	}
	recent = t
	return
}

func (s *CatalogRepository) getReservedCounts(ctx context.Context, itemID string) (count int, err error) {
	reserved := r.Table(tableMarket).
		GetAllByIndex(marketFieldItemID, itemID).
		Filter(dotagiftx.Market{
			Type:   dotagiftx.MarketTypeAsk,
			Status: dotagiftx.MarketStatusReserved,
		}).Count()

	if err = s.db.one(ctx, reserved, &count); err != nil {
		err = fmt.Errorf("could not get reserved count: %s", err)
	}

	return
}

func (s *CatalogRepository) medianPriceQuery(qty int, t r.Term) r.Term {
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

func (s *CatalogRepository) create(ctx context.Context, in *dotagiftx.Catalog) error {
	// Fixes missing item in catalog that does not have views yet.
	in.ViewCount = 1
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	// Convert catalog into map to insert zero value fields.
	m := catalogToMap(in)

	if _, err := s.db.insert(ctx, s.table().Insert(m)); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return nil
}

func (s *CatalogRepository) update(ctx context.Context, in *dotagiftx.Catalog) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	// Convert catalog into map to insert zero value fields.
	m := catalogToMap(in)

	err = s.db.update(ctx, s.table().Get(in.ID).Update(m))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *CatalogRepository) table() r.Term {
	return r.Table(tableCatalog)
}

func catalogToMap(cat *dotagiftx.Catalog) map[string]any {
	s := structs.New(cat)
	s.TagName = "json"
	return s.Map()
}
