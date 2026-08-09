package rethink

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"dario.cat/mergo"
	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableUser              = "user"
	userFieldSteamID       = "steam_id"
	userSubscriptionEndsAt = "subscription_ends_at"
)

var userSearchFields = []string{"name", "steam_id", "url"}

// NewUser creates new instance of user data store.
func NewUser(c *Client) *UserRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableUser); err != nil {
		log.Fatalf("could not create %s table: %s", tableUser, err)
	}

	if err := c.autoIndex(ctx, tableUser, dotagiftx2.User{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableUser, err)
	}

	return &UserRepository{c, userSearchFields}
}

type UserRepository struct {
	db            *Client
	keywordFields []string
}

func (s *UserRepository) Find(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.User, error) {
	var res []dotagiftx2.User
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *UserRepository) FindFlagged(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.User, error) {
	var res []dotagiftx2.User
	o.KeywordFields = s.keywordFields
	q := baseFindOptsQuery(s.table(), o, s.flaggedFilter)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *UserRepository) flaggedFilter(q r.Term) r.Term {
	return q.Filter(func(t r.Term) any {
		return t.Field("status").Ge(dotagiftx2.UserStatusSuspended)
	})
}

func (s *UserRepository) Count(ctx context.Context, o dotagiftx2.FindOpts) (num int, err error) {
	o = dotagiftx2.FindOpts{Filter: o.Filter, UserID: o.UserID}
	q := newFindOptsQuery(s.table(), o)
	err = s.db.one(ctx, q.Count(), &num)
	return
}

func (s *UserRepository) Get(ctx context.Context, id string) (*dotagiftx2.User, error) {
	// Check steam ID first exist.
	row, _ := s.getBySteamID(ctx, id)
	if row != nil {
		return row, nil
	}

	// Try to find it by user ID.
	row = &dotagiftx2.User{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.UserErrNotFound
		}

		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *UserRepository) getBySteamID(ctx context.Context, steamID string) (*dotagiftx2.User, error) {
	row := &dotagiftx2.User{}
	q := s.table().GetAllByIndex(userFieldSteamID, steamID).OrderBy("created_at").Limit(1)
	if err := s.db.one(ctx, q, row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.UserErrNotFound
		}

		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *UserRepository) Create(ctx context.Context, in *dotagiftx2.User) error {
	t := now()
	if in.CreatedAt == nil {
		in.CreatedAt = t
	}
	in.UpdatedAt = t
	id, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *UserRepository) Update(ctx context.Context, in *dotagiftx2.User) error {
	in.UpdatedAt = now()
	return s.BaseUpdate(ctx, in)
}

func (s *UserRepository) BaseUpdate(ctx context.Context, in *dotagiftx2.User) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	err = s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	if err = mergo.Merge(in, cur); err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageMergeErr, err)
	}

	return nil
}

// ExpiringSubscribers return expiring subscribers on given t time.
func (s *UserRepository) ExpiringSubscribers(ctx context.Context, t time.Time) ([]dotagiftx2.User, error) {
	var res []dotagiftx2.User
	q := s.table().HasFields(userSubscriptionEndsAt)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	var expiring []dotagiftx2.User
	for _, u := range res {
		if u.SubscriptionEndsAt.After(t) {
			continue
		}
		expiring = append(expiring, u)
	}
	return expiring, nil
}

// PurgeSubscription clears subscription data.
func (s *UserRepository) PurgeSubscription(ctx context.Context, userID string) error {
	t := time.Now()
	err := s.db.update(ctx, s.table().Get(userID).Update(map[string]any{
		"boons":                r.Literal(),
		"subscription":         r.Literal(),
		"subscribed_at":        r.Literal(),
		userSubscriptionEndsAt: r.Literal(),
		"subscription_notes":   fmt.Sprintf("purged at %s", t),
		"updated_at":           t,
	}))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}
	return nil
}

func (s *UserRepository) ClearSubscriptionEndsAt(ctx context.Context, userID string) error {
	q := s.table().Get(userID).Update(map[string]any{
		userSubscriptionEndsAt: r.Literal(),
	})
	if err := s.db.update(ctx, q); err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return nil
}

func (s *UserRepository) table() r.Term {
	return r.Table(tableUser)
}
