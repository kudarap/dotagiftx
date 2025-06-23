package rethink

import (
	"fmt"
	"log"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableMarket                = "market"
	marketFieldItemID          = "item_id"
	marketFieldUserID          = "user_id"
	marketFieldType            = "type"
	marketFieldStatus          = "status"
	marketFieldInventoryStatus = "inventory_status"
	marketFieldDeliveryStatus  = "delivery_status"
	marketFieldNotes           = "notes"
	marketFieldPrice           = "price"
	marketFieldResell          = "resell"
	marketFieldCreatedAt       = "created_at"
	marketFieldUpdatedAt       = "updated_at"
	// Hidden field for searching item details.
	marketItemSearchTags = "search_text"
)

// NewMarket creates new instance of market data store.
func NewMarket(c *Client) dotagiftx.MarketStorage {
	if err := c.autoMigrate(tableMarket); err != nil {
		log.Fatalf("could not create %s table: %s", tableMarket, err)
	}
	if err := c.autoIndex(tableMarket, dotagiftx.Market{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &marketStorage{c, []string{marketItemSearchTags}}
}

type marketStorage struct {
	db            *Client
	keywordFields []string
}

func (s *marketStorage) Find(o dotagiftx.FindOpts) ([]dotagiftx.Market, error) {
	var res []dotagiftx.Market
	o.KeywordFields = s.keywordFields
	o.IndexSorting = true

	q := findOpts(o).parseOpts(s.table(), nil)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	for i, rr := range res {
		res[i].User = s.includeUser(rr.UserID)
	}

	return res, nil
}

// PendingInventoryStatus returns market entries that is pending for checking
// inventory status or needs re-processing of re-process error status.
func (s *marketStorage) PendingInventoryStatus(o dotagiftx.FindOpts) ([]dotagiftx.Market, error) {
	// Filters out already check or no need to check market
	q := newFindOptsQuery(r.Table(tableMarket), o)
	q = q.
		// .filter(r.row.hasFields('inventory_status').not().or(r.row('inventory_status').eq(500)))
		Filter(func(t r.Term) r.Term {
			return t.HasFields(marketFieldInventoryStatus).Not().
				Or(t.Field(marketFieldInventoryStatus).Eq(dotagiftx.InventoryStatusError))
		}).
		Filter(func(t r.Term) r.Term {
			return t.And(t.Field(marketFieldStatus).Eq(dotagiftx.MarketStatusLive).
				Or(t.Field(marketFieldStatus).Eq(dotagiftx.MarketStatusReserved))).
				And(t.Field(marketFieldType).Eq(dotagiftx.MarketTypeAsk))
		})

	var res []dotagiftx.Market
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	for i, rr := range res {
		res[i].User = s.includeUser(rr.UserID)
	}
	return res, nil
}

// PendingDeliveryStatus returns market entries that is pending for checking
// delivery status or needs re-processing of re-process error status.
func (s *marketStorage) PendingDeliveryStatus(o dotagiftx.FindOpts) ([]dotagiftx.Market, error) {
	q := newFindOptsQuery(r.Table(tableMarket), o)
	q = q.Filter(func(t r.Term) r.Term {
		return t.HasFields(marketFieldDeliveryStatus).Not().
			Or(t.Field(marketFieldDeliveryStatus).Eq(dotagiftx.DeliveryStatusError))
		//Or(t.Field(marketFieldDeliveryStatus).Eq(core.DeliveryStatusError).
		//	Or(t.Field(marketFieldDeliveryStatus).Eq(core.DeliveryStatusNoHit)))
	})

	var res []dotagiftx.Market
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	for i, rr := range res {
		res[i].User = s.includeUser(rr.UserID)
	}
	return res, nil
}

func (s *marketStorage) RevalidateDeliveryStatus(o dotagiftx.FindOpts) ([]dotagiftx.Market, error) {
	n := time.Now()
	q := newFindOptsQuery(r.Table(tableMarket), o)
	q = q.Filter(func(t r.Term) r.Term {
		return t.Field(marketFieldStatus).Eq(dotagiftx.MarketStatusSold).
			And(t.Field(marketFieldUpdatedAt).Year().Eq(n.Year()).
				And(t.Field(marketFieldUpdatedAt).Month().Eq(n.Month()).
					And(t.Field(marketFieldUpdatedAt).Day().Eq(n.Day())))).
			And(t.Field(marketFieldDeliveryStatus).Eq(dotagiftx.DeliveryStatusNoHit).
				Or(t.Field(marketFieldDeliveryStatus).Eq(dotagiftx.DeliveryStatusPrivate)))
	})

	var res []dotagiftx.Market
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	for i, rr := range res {
		res[i].User = s.includeUser(rr.UserID)
	}
	return res, nil
}

func (s *marketStorage) Count(o dotagiftx.FindOpts) (num int, err error) {
	o = dotagiftx.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
		IndexKey:      o.IndexKey,
	}

	q := findOpts(o).parseOpts(s.table(), nil)
	err = s.db.one(q.Count(), &num)
	return
}

// includeRelatedFields injects item and user details base on market foreign keys
// and create a search tag
func (s *marketStorage) includeRelatedFields(q r.Term) r.Term {
	return q.
		//EqJoin(marketFieldItemID, r.Table(tableItem)).
		//Map(func(t r.Term) r.Term {
		//	market := t.Field("left")
		//	item := t.Field("right")
		//	tags := market.Field(marketFieldNotes).Default("")
		//	for _, ff := range itemSearchFields {
		//		tags = tags.Add(" ", item.Field(ff))
		//	}
		//
		//	return market.Merge(map[string]interface{}{
		//		tableItem:            item,
		//		marketItemSearchTags: tags,
		//	})
		//}).
		EqJoin(marketFieldUserID, r.Table(tableUser)).
		Map(func(t r.Term) r.Term {
			return t.Field("left").Merge(map[string]interface{}{
				tableUser: t.Field("right"),
			})
		})
}

func (s *marketStorage) Get(id string) (*dotagiftx.Market, error) {
	row := &dotagiftx.Market{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dotagiftx.MarketErrNotFound
		}

		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	row.User = s.includeUser(row.UserID)
	return row, nil
}

func (s *marketStorage) includeUser(userID string) *dotagiftx.User {
	var user dotagiftx.User
	_ = s.db.one(r.Table(tableUser).Get(userID), &user)
	return &user
}

func (s *marketStorage) Index(id string) (*dotagiftx.Market, error) {
	mkt, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	var item dotagiftx.Item
	_ = s.db.one(r.Table(tableItem).Get(mkt.ItemID), &item)
	mkt.Item = &item

	var invs []dotagiftx.Inventory
	_ = s.db.list(r.Table(tableInventory).GetAllByIndex(inventoryFieldMarketID, mkt.ID), &invs)
	if len(invs) != 0 {
		mkt.Inventory = &invs[0]
	}

	var dels []dotagiftx.Delivery
	_ = s.db.list(r.Table(tableDelivery).GetAllByIndex(deliveryFieldMarketID, mkt.ID), &dels)
	if len(dels) != 0 {
		mkt.Delivery = &dels[0]
	}

	mkt.SearchText = mkt.Notes
	if mkt.Item != nil {
		mkt.SearchText += strings.Join([]string{
			"",
			mkt.Item.Name,
			mkt.Item.Hero,
			mkt.Item.Origin,
			mkt.Item.Rarity,
		}, " ")
	}

	if err = s.BaseUpdate(mkt); err != nil {
		return nil, err
	}

	return mkt, nil
}

func (s *marketStorage) Create(in *dotagiftx.Market) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	in.User = nil
	in.Item = nil
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *marketStorage) Update(in *dotagiftx.Market) error {
	in.UpdatedAt = now()
	return s.BaseUpdate(in)
}

func (s *marketStorage) UpdateUserScore(userID string, rankScore int) error {
	if userID == "" {
		return fmt.Errorf("user id is required to update user score")
	}

	// get all user live market
	var markets []dotagiftx.Market
	q := s.table().GetAllByIndex(marketFieldUserID, userID).Filter(dotagiftx.Market{Status: dotagiftx.MarketStatusLive})
	if err := s.db.list(q, &markets); err != nil {
		return err
	}

	// set new user rank score
	for _, mm := range markets {
		mm.UserRankScore = rankScore
		if err := s.BaseUpdate(&mm); err != nil {
			return fmt.Errorf("could not update market user rank: %s", err)
		}
	}

	return nil
}

func (s *marketStorage) BaseUpdate(in *dotagiftx.Market) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.User = nil
	//in.Item = nil
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return errors.New(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *marketStorage) UpdateExpiring(t dotagiftx.MarketType, b dotagiftx.UserBoon, cutOff time.Time) (ids []string, err error) {
	// Collects exempted users ids.
	q := r.Table(tableUser).
		HasFields("boons").
		Filter(r.Row.Field("boons").Contains(b)).
		Field("id")
	var exemptedUserIDs []string
	if err = s.db.list(q, &exemptedUserIDs); err != nil {
		return nil, fmt.Errorf("could not get users: %s", err)
	}

	// Sets expired entry state base on cutOff time.
	now := time.Now()
	q = s.table().GetAllByIndex("status", dotagiftx.MarketStatusLive).Filter(dotagiftx.Market{Type: t}).
		Filter(r.Row.Field(marketFieldCreatedAt).Lt(cutOff)).
		Filter(func(entry r.Term) r.Term {
			return r.Expr(exemptedUserIDs).Contains(entry.Field(marketFieldUserID)).Not()
		}).
		Update(dotagiftx.Market{
			Status: dotagiftx.MarketStatusExpired, UpdatedAt: &now,
		})
	if err = s.db.update(q); err != nil {
		return nil, fmt.Errorf("could not update expiring markets: %s", err)
	}

	// Collect and return affected item ids.
	q = s.table().GetAllByIndex("status", dotagiftx.MarketStatusExpired).Filter(dotagiftx.Market{UpdatedAt: &now}).
		Group(marketFieldItemID).Count().Ungroup().
		Field("group")
	var itemIDs []string
	if err = s.db.list(q, &itemIDs); err != nil {
		return nil, fmt.Errorf("could not get affected markets: %s", err)
	}

	return itemIDs, nil
}

func (s *marketStorage) UpdateExpiringResell(b dotagiftx.UserBoon) (ids []string, err error) {
	// Collects exempted users ids.
	q := r.Table(tableUser).
		HasFields("boons").
		Filter(r.Row.Field("boons").Contains(b)).
		Field("id")
	var exemptedUserIDs []string
	if err = s.db.list(q, &exemptedUserIDs); err != nil {
		return nil, fmt.Errorf("could not get users: %s", err)
	}

	// Sets expired entry state immediately.
	now := time.Now()
	resell := true
	q = s.table().GetAllByIndex("status", dotagiftx.MarketStatusLive).
		Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeAsk, Resell: &resell}).
		Filter(func(entry r.Term) r.Term {
			return r.Expr(exemptedUserIDs).Contains(entry.Field(marketFieldUserID)).Not()
		}).
		Update(dotagiftx.Market{
			Status: dotagiftx.MarketStatusExpired, UpdatedAt: &now,
		})
	if err = s.db.update(q); err != nil {
		return nil, fmt.Errorf("could not update expiring markets: %s", err)
	}

	// Collect and return affected item ids.
	q = s.table().GetAllByIndex("status", dotagiftx.MarketStatusExpired).Filter(dotagiftx.Market{UpdatedAt: &now}).
		Group(marketFieldItemID).Count().Ungroup().
		Field("group")
	var itemIDs []string
	if err = s.db.list(q, &itemIDs); err != nil {
		return nil, fmt.Errorf("could not get affected markets: %s", err)
	}

	return itemIDs, nil
}

func (s *marketStorage) BulkDeleteByStatus(ms dotagiftx.MarketStatus, cutOff time.Time, limit int) error {
	if ms != dotagiftx.MarketStatusRemoved && ms != dotagiftx.MarketStatusExpired {
		return fmt.Errorf("market status %s not allowed to bulk delete", ms)
	}

	q := s.table().GetAllByIndex("status", ms).
		Filter(r.Row.Field(marketFieldCreatedAt).Lt(cutOff)).
		Limit(limit).
		Delete()
	if err := s.db.delete(q); err != nil && err != r.ErrEmptyResult {
		return err
	}
	return nil
}

func (s *marketStorage) findIndexLegacy(o dotagiftx.FindOpts) ([]dotagiftx.Catalog, error) {
	q := s.indexBaseQuery()

	var res []dotagiftx.Catalog
	o.KeywordFields = s.keywordFields
	q = newFindOptsQuery(q, o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *marketStorage) countIndexLegacy(o dotagiftx.FindOpts) (num int, err error) {
	q := s.indexBaseQuery()
	o = dotagiftx.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
	}
	q = newFindOptsQuery(q, o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *marketStorage) indexBaseQuery() r.Term {
	return s.table().GroupByIndex(marketFieldItemID).Ungroup().
		Map(s.groupIndexMap).
		EqJoin(marketFieldItemID, r.Table(tableItem)).
		Zip()
}

func (s *marketStorage) table() r.Term {
	return r.Table(tableMarket)
}

func (s *marketStorage) groupIndexMap(market r.Term) interface{} {
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

	id := market.Field("group")
	live := market.Field("reduction").Filter(dotagiftx.Market{Status: dotagiftx.MarketStatusLive})
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
