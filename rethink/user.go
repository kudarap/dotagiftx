package rethink

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableUser        = "user"
	itemFieldSteamID = "steam_id"
)

// NewUser creates new instance of user data store.
func NewUser(c *Client) core.UserStorage {
	if err := c.autoMigrate(tableUser); err != nil {
		log.Fatalf("could not create %s table: %s", tableUser, err)
	}

	if err := c.createIndex(tableUser, itemFieldSteamID); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableUser, err)
	}

	return &userStorage{c}
}

type userStorage struct {
	db *Client
}

func (s *userStorage) Find(o core.FindOpts) ([]core.User, error) {
	var res []core.User
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *userStorage) Count(o core.FindOpts) (num int, err error) {
	o = core.FindOpts{Filter: o.Filter, UserID: o.UserID}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *userStorage) Get(id string) (*core.User, error) {
	// Check steam ID first exist.
	row, _ := s.getBySteamID(id)
	if row != nil {
		return row, nil
	}

	// Try find it by user ID.
	row = &core.User{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.UserErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *userStorage) getBySteamID(steamID string) (*core.User, error) {
	row := &core.User{}
	q := s.table().GetAllByIndex(itemFieldSteamID, steamID)
	if err := s.db.one(q, row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.UserErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *userStorage) Create(in *core.User) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *userStorage) Update(in *core.User) error {
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

func (s *userStorage) table() r.Term {
	return r.Table(tableUser)
}
