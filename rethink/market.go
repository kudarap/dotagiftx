package rethink

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableMarket          = "market"
	marketFieldItemID    = "item_id"
	marketFieldUserID    = "user_id"
	marketFieldType      = "type"
	marketFieldStatus    = "status"
	marketFieldNotes     = "notes"
	marketFieldPrice     = "price"
	marketFieldCreatedAt = "created_at"
	marketFieldUpdatedAt = "updated_at"
	// Hidden field for searching item details.
	marketItemSearchTags = "item_tags"
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
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
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
	q := findOpts(o).parseOpts(s.table(), s.includeRelatedFields)
	err = s.db.one(q.Count(), &num)
	return
}

// includeRelatedFields injects item and user details base on market foreign keys
// and create a search tag
func (s *marketStorage) includeRelatedFields(q r.Term) r.Term {
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
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	in.User = nil
	in.Item = nil
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
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
