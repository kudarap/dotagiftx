package rethink

import (
	"log"
	"time"

	"github.com/imdario/mergo"
	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
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

func (s *catalogStorage) Find(o core.FindOpts) ([]core.Catalog, error) {
	var res []core.Catalog
	o.KeywordFields = s.keywordFields
	o.IndexSorting = true
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
		IndexSorting:  true,
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
	// Benchmark indexing.
	// avg proc time 2.555710726s
	tStart := time.Now()
	defer func() {
		s.logger.Infof("catalog indexed %s @ %s\n", itemID, time.Now().Sub(tStart))
	}()

	cat := &core.Catalog{}
	opts := core.FindOpts{Filter: core.Market{ItemID: itemID}}
	baseQ := newFindOptsQuery(r.Table(tableMarket), opts)

	var q r.Term
	var err error

	// Get item details by item ID.
	q = r.Table(tableItem).Get(itemID)
	if err = s.db.one(q, cat); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get total market count by item ID.
	if err = s.db.one(baseQ.Count(), &cat.Quantity); err != nil {
		return nil, errors.New(core.CatalogErrIndexing, err)
	}

	// Get lowest price on the market by item ID.
	q = baseQ.Min("price").Field("price").Default(0)
	if err = s.db.one(q, &cat.LowestAsk); err != nil {
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
