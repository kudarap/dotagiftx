package rethink

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tableCatalog = "catalog"

// NewCatalog creates new instance of catalog data store.
func NewCatalog(c *Client) core.CatalogStorage {
	if err := c.autoMigrate(tableCatalog); err != nil {
		log.Fatalf("could not create %s table: %s", tableCatalog, err)
	}

	if err := c.autoIndex(tableCatalog, core.Catalog{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableCatalog, err)
	}

	return &catalogStorage{c, []string{"name", "hero", "origin", "rarity"}}
}

type catalogStorage struct {
	db            *Client
	keywordFields []string
}

func (s *catalogStorage) Find(o core.FindOpts) ([]core.Catalog, error) {
	var res []core.Catalog
	o.KeywordFields = s.keywordFields
	q := newFindOptsQuery(s.table(), o)
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
		UserID:        o.UserID,
	}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *catalogStorage) Get(itemID string) (*core.Catalog, error) {
	row := &core.Catalog{}
	if err := s.db.one(s.table().Get(itemID), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.CatalogErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *catalogStorage) Index(itemID string) (*core.Catalog, error) {
	row := &core.Catalog{}

	// Get market index base item ID.

	// Check for exiting entry

	return row, nil
}

func (s *catalogStorage) create(in *core.Catalog) error {
	t := now()
	in.UpdatedAt = t
	in.ID = ""
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}
	in.ID = id

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
