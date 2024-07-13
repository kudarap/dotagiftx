package rethink

import (
	"context"
	"fmt"
	"log"
	"time"

	"dario.cat/mergo"
	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableUser        = "user"
	userFieldSteamID = "steam_id"
)

var userSearchFields = []string{"name", "steam_id", "url"}

// NewUser creates new instance of user data store.
func NewUser(c *Client) dgx.UserStorage {
	if err := c.autoMigrate(tableUser); err != nil {
		log.Fatalf("could not create %s table: %s", tableUser, err)
	}

	if err := c.autoIndex(tableUser, dgx.User{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableUser, err)
	}

	return &userStorage{c, userSearchFields}
}

type userStorage struct {
	db            *Client
	keywordFields []string
}

func (s *userStorage) Find(o dgx.FindOpts) ([]dgx.User, error) {
	var res []dgx.User
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *userStorage) FindFlagged(o dgx.FindOpts) ([]dgx.User, error) {
	var res []dgx.User
	o.KeywordFields = s.keywordFields
	q := baseFindOptsQuery(s.table(), o, s.flaggedFilter)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *userStorage) flaggedFilter(q r.Term) r.Term {
	return q.Filter(func(t r.Term) interface{} {
		return t.Field("status").Ge(dgx.UserStatusSuspended)
	})
}

func (s *userStorage) Count(o dgx.FindOpts) (num int, err error) {
	o = dgx.FindOpts{Filter: o.Filter, UserID: o.UserID}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(q.Count(), &num)
	return
}

func (s *userStorage) Get(id string) (*dgx.User, error) {
	// Check steam ID first exist.
	row, _ := s.getBySteamID(id)
	if row != nil {
		return row, nil
	}

	// Try find it by user ID.
	row = &dgx.User{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.UserErrNotFound
		}

		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *userStorage) getBySteamID(steamID string) (*dgx.User, error) {
	row := &dgx.User{}
	q := s.table().GetAllByIndex(userFieldSteamID, steamID)
	if err := s.db.one(q, row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.UserErrNotFound
		}

		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *userStorage) Create(in *dgx.User) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *userStorage) Update(in *dgx.User) error {
	in.UpdatedAt = now()
	return s.BaseUpdate(in)
}

func (s *userStorage) BaseUpdate(in *dgx.User) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return errors.New(dgx.StorageMergeErr, err)
	}

	return nil
}

// ExpiringSubscribers return expiring subscribers on given t time.
func (s *userStorage) ExpiringSubscribers(ctx context.Context, t time.Time) ([]dgx.User, error) {
	var res []dgx.User
	q := s.table().HasFields("subscription_ends_at")
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	var expiring []dgx.User
	for _, u := range res {
		if u.SubscriptionEndsAt.After(t) {
			continue
		}
		expiring = append(expiring, u)
	}
	return expiring, nil
}

// PurgeSubscription clears subscription data.
func (s *userStorage) PurgeSubscription(ctx context.Context, userID string) error {
	t := time.Now()
	err := s.db.update(s.table().Get(userID).Update(map[string]interface{}{
		"boons":                r.Literal(),
		"subscription":         r.Literal(),
		"subscribed_at":        r.Literal(),
		"subscription_ends_at": r.Literal(),
		"subscription_notes":   fmt.Sprintf("purged at %s", t),
		"updated_at":           t,
	}))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}
	return nil
}

func (s *userStorage) table() r.Term {
	return r.Table(tableUser)
}
