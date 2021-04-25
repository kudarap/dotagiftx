package rethink

import (
	"log"
	"strings"

	"github.com/imdario/mergo"
	"github.com/kudarap/dotagiftx/core"
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
	marketFieldCreatedAt       = "created_at"
	marketFieldUpdatedAt       = "updated_at"
	// Hidden field for searching item details.
	marketItemSearchTags = "search_text"
)

// NewMarket creates new instance of market data store.
func NewMarket(c *Client) core.MarketStorage {
	if err := c.autoMigrate(tableMarket); err != nil {
		log.Fatalf("could not create %s table: %s", tableMarket, err)
	}
	if err := c.autoIndex(tableMarket, core.Market{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &marketStorage{c, []string{marketItemSearchTags}}
}

type marketStorage struct {
	db            *Client
	keywordFields []string
}

func (s *marketStorage) Find(o core.FindOpts) ([]core.Market, error) {
	var res []core.Market
	o.KeywordFields = s.keywordFields
	// IndexSorting was disable due to hook query includeRelatedFields
	//o.IndexSorting = true
	q := findOpts(o).parseOpts(s.table(), nil)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

// PendingInventoryStatus returns market entries that is pending for checking
// inventory status or needs re-processing of re-process error status.
func (s *marketStorage) PendingInventoryStatus(o core.FindOpts) ([]core.Market, error) {
	// Filters out already check or no need to check market
	q := r.Table(tableMarket).
		// .filter(r.row.hasFields('inventory_status').not().or(r.row('inventory_status').eq(500)))
		Filter(func(t r.Term) r.Term {
			return t.HasFields(marketFieldInventoryStatus).Not().Or(
				t.Field(marketFieldInventoryStatus).Eq(core.InventoryStatusError),
			)
		})
	q = baseFindOptsQuery(q, o, s.includeRelatedFields)

	var res []core.Market
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

// PendingDeliveryStatus returns market entries that is pending for checking
// delivery status or needs re-processing of re-process error status.
func (s *marketStorage) PendingDeliveryStatus(o core.FindOpts) ([]core.Market, error) {
	q := r.Table(tableMarket).
		Filter(func(t r.Term) r.Term {
			return t.HasFields(marketFieldDeliveryStatus).Not().
				Or(t.Field(marketFieldDeliveryStatus).Eq(core.DeliveryStatusError))
			//Or(t.Field(marketFieldDeliveryStatus).Eq(core.DeliveryStatusError).
			//	Or(t.Field(marketFieldDeliveryStatus).Eq(core.DeliveryStatusNoHit)))
		})
	q = baseFindOptsQuery(q, o, s.includeRelatedFields)

	var res []core.Market
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *marketStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{
		Keyword:       o.Keyword,
		KeywordFields: s.keywordFields,
		Filter:        o.Filter,
		UserID:        o.UserID,
		// IndexSorting was disable due to hook query includeRelatedFields
		//IndexSorting:  true,
	}
	q := findOpts(o).parseOpts(s.table(), nil)
	err = s.db.one(q.Count(), &num)
	return
}

// includeRelatedFields injects item and user details base on market foreign keys
// and create a search tag
func (s *marketStorage) includeRelatedFields(q r.Term) r.Term {
	return q
}

// slowIncludeRelatedFields deprecated
func (s *marketStorage) slowIncludeRelatedFields(q r.Term) r.Term {
	return q.
		EqJoin(marketFieldItemID, r.Table(tableItem)).
		Map(func(t r.Term) r.Term {
			market := t.Field("left")
			item := t.Field("right")
			tags := market.Field(marketFieldNotes).Default("")
			for _, ff := range itemSearchFields {
				tags = tags.Add(" ", item.Field(ff))
			}

			return market.Merge(map[string]interface{}{
				tableItem:            item,
				marketItemSearchTags: tags,
			})
		}).
		EqJoin(marketFieldUserID, r.Table(tableUser)).
		Map(func(t r.Term) r.Term {
			return t.Field("left").Merge(map[string]interface{}{
				tableUser: t.Field("right"),
			})
		})
}

func (s *marketStorage) Get(id string) (*core.Market, error) {
	row := &core.Market{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.MarketErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *marketStorage) Index(id string) (*core.Market, error) {
	mkt, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	_ = s.db.one(r.Table(tableItem).Get(mkt.ItemID), mkt.Item)

	var invs []core.Inventory
	_ = s.db.list(r.Table(tableInventory).GetAllByIndex(inventoryFieldMarketID, mkt.ID), &invs)
	if len(invs) != 0 {
		mkt.Inventory = &invs[0]
	}

	var dels []core.Delivery
	_ = s.db.list(r.Table(tableDelivery).GetAllByIndex(inventoryFieldMarketID, mkt.ID), &dels)
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

func (s *marketStorage) Create(in *core.Market) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	in.ID = ""
	in.User = nil
	in.Item = nil
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *marketStorage) Update(in *core.Market) error {
	in.UpdatedAt = now()
	return s.BaseUpdate(in)
}

func (s *marketStorage) BaseUpdate(in *core.Market) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.User = nil
	in.Item = nil
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return errors.New(core.StorageMergeErr, err)
	}

	return nil
}

func (s *marketStorage) findIndexLegacy(o core.FindOpts) ([]core.Catalog, error) {
	q := s.indexBaseQuery()

	var res []core.Catalog
	o.KeywordFields = s.keywordFields
	q = newFindOptsQuery(q, o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *marketStorage) countIndexLegacy(o core.FindOpts) (num int, err error) {
	q := s.indexBaseQuery()
	o = core.FindOpts{
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
	live := market.Field("reduction").Filter(core.Market{Status: core.MarketStatusLive})
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
