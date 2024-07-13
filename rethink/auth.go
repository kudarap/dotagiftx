package rethink

import (
	"log"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"

	"dario.cat/mergo"
)

const (
	tableAuth             = "auth"
	authFieldUsername     = "username"
	authFieldRefreshToken = "refresh_token"
)

// NewAuth creates new instance of auth data store.
func NewAuth(c *Client) *authStorage {
	if err := c.autoMigrate(tableAuth); err != nil {
		log.Fatalf("could not create %s table: %s", tableAuth, err)
	}

	if err := c.autoIndex(tableAuth, dgx.Auth{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableAuth, err)
	}

	return &authStorage{c}
}

type authStorage struct {
	db *Client
}

func (s *authStorage) Get(id string) (*dgx.Auth, error) {
	row := &dgx.Auth{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.AuthErrNotFound
		}

		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *authStorage) GetByUsername(username string) (*dgx.Auth, error) {
	row := &dgx.Auth{}
	q := s.table().GetAllByIndex(authFieldUsername, username)
	if err := s.db.one(q, row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dgx.AuthErrNotFound
		}

		return nil, err
	}

	return row, nil
}

func (s *authStorage) GetByUsernameAndPassword(username, password string) (*dgx.Auth, error) {
	return s.findOne(dgx.Auth{Username: username, Password: password})
}

func (s *authStorage) GetByRefreshToken(refreshToken string) (*dgx.Auth, error) {
	row := &dgx.Auth{}
	q := s.table().GetAllByIndex(authFieldRefreshToken, refreshToken)
	if err := s.db.one(q, row); err != nil {
		return nil, err
	}

	return row, nil
}

func (s *authStorage) Create(in *dgx.Auth) error {
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

func (s *authStorage) Update(in *dgx.Auth) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(dgx.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return errors.New(dgx.StorageMergeErr, err)
	}

	return nil
}

func (s *authStorage) find(o dgx.FindOpts) ([]dgx.Auth, error) {
	var res []dgx.Auth
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(dgx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *authStorage) findOne(filter dgx.Auth) (*dgx.Auth, error) {
	o := dgx.FindOpts{Filter: filter, Limit: 1}
	res, err := s.find(o)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, dgx.AuthErrNotFound
	}

	return &res[0], nil
}

func (s *authStorage) table() r.Term {
	return r.Table(tableAuth)
}
