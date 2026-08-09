package rethink

import (
	"context"
	"errors"
	"log"

	"dario.cat/mergo"
	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
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

	if err := c.autoIndex(ctx, tableAuth, dotagiftx2.Auth{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableAuth, err)
	}

	return &AuthRepository{c}
}

type AuthRepository struct {
	db *Client
}

func (s *AuthRepository) Get(ctx context.Context, id string) (*dotagiftx2.Auth, error) {
	row := &dotagiftx2.Auth{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.AuthErrNotFound
		}

		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *AuthRepository) GetByUsername(ctx context.Context, username string) (*dotagiftx2.Auth, error) {
	row := &dotagiftx2.Auth{}
	q := s.table().GetAllByIndex(authFieldUsername, username).OrderBy("created_at").Limit(1)
	if err := s.db.one(ctx, q, row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx2.AuthErrNotFound
		}

		return nil, err
	}

	return row, nil
}

func (s *AuthRepository) GetByUsernameAndPassword(ctx context.Context, username, password string) (*dotagiftx2.Auth, error) {
	return s.findOne(ctx, dotagiftx2.Auth{Username: username, Password: password})
}

func (s *AuthRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*dotagiftx2.Auth, error) {
	row := &dotagiftx2.Auth{}
	q := s.table().GetAllByIndex(authFieldRefreshToken, refreshToken)
	if err := s.db.one(ctx, q, row); err != nil {
		return nil, err
	}

	return row, nil
}

func (s *AuthRepository) Create(ctx context.Context, in *dotagiftx2.Auth) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(ctx, s.table().Insert(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *AuthRepository) Update(ctx context.Context, in *dotagiftx2.Auth) error {
	cur, err := s.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageMergeErr, err)
	}

	return nil
}

func (s *AuthRepository) find(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Auth, error) {
	var res []dotagiftx2.Auth
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *AuthRepository) findOne(ctx context.Context, filter dotagiftx2.Auth) (*dotagiftx2.Auth, error) {
	o := dotagiftx2.FindOpts{Filter: filter, Limit: 1}
	res, err := s.find(ctx, o)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, dotagiftx2.AuthErrNotFound
	}

	return &res[0], nil
}

func (s *AuthRepository) table() r.Term {
	return r.Table(tableAuth)
}
