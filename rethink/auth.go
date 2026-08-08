package rethink

import (
	"context"
	"errors"
	"log"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableAuth             = "auth"
	authFieldUsername     = "username"
	authFieldRefreshToken = "refresh_token"
)

// NewAuth creates a new instance of auth data store.
func NewAuth(c *Client) *AuthRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableAuth); err != nil {
		log.Fatalf("could not create %s table: %s", tableAuth, err)
	}

	if err := c.autoIndex(ctx, tableAuth, dotagiftx.Auth{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableAuth, err)
	}

	return &AuthRepository{c}
}

type AuthRepository struct {
	db *Client
}

func (s *AuthRepository) Get(ctx context.Context, id string) (*dotagiftx.Auth, error) {
	row := &dotagiftx.Auth{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.AuthErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *AuthRepository) GetByUsername(ctx context.Context, username string) (*dotagiftx.Auth, error) {
	row := &dotagiftx.Auth{}
	q := s.table().GetAllByIndex(authFieldUsername, username).OrderBy("created_at").Limit(1)
	if err := s.db.one(ctx, q, row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.AuthErrNotFound
		}

		return nil, err
	}

	return row, nil
}

func (s *AuthRepository) GetByUsernameAndPassword(ctx context.Context, username, password string) (*dotagiftx.Auth, error) {
	return s.findOne(ctx, dotagiftx.Auth{Username: username, Password: password})
}

func (s *AuthRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*dotagiftx.Auth, error) {
	row := &dotagiftx.Auth{}
	q := s.table().GetAllByIndex(authFieldRefreshToken, refreshToken)
	if err := s.db.one(ctx, q, row); err != nil {
		return nil, err
	}

	return row, nil
}

func (s *AuthRepository) Create(ctx context.Context, in *dotagiftx.Auth) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *AuthRepository) Update(ctx context.Context, in *dotagiftx.Auth) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *AuthRepository) find(ctx context.Context, o dotagiftx.FindOpts) ([]dotagiftx.Auth, error) {
	var res []dotagiftx.Auth
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *AuthRepository) findOne(ctx context.Context, filter dotagiftx.Auth) (*dotagiftx.Auth, error) {
	o := dotagiftx.FindOpts{Filter: filter, Limit: 1}
	res, err := s.find(ctx, o)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, dotagiftx.AuthErrNotFound
	}

	return &res[0], nil
}

func (s *AuthRepository) table() r.Term {
	return r.Table(tableAuth)
}
