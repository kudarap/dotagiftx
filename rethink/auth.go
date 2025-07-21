package rethink

import (
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
func NewAuth(c *Client) *authStorage {
	if err := c.autoMigrate(tableAuth); err != nil {
		log.Fatalf("could not create %s table: %s", tableAuth, err)
	}

	if err := c.autoIndex(tableAuth, dotagiftx.Auth{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableAuth, err)
	}

	return &authStorage{c}
}

type authStorage struct {
	db *Client
}

func (s *authStorage) Get(id string) (*dotagiftx.Auth, error) {
	row := &dotagiftx.Auth{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dotagiftx.AuthErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *authStorage) GetByUsername(username string) (*dotagiftx.Auth, error) {
	row := &dotagiftx.Auth{}
	q := s.table().GetAllByIndex(authFieldUsername, username)
	if err := s.db.one(q, row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, dotagiftx.AuthErrNotFound
		}

		return nil, err
	}

	return row, nil
}

func (s *authStorage) GetByUsernameAndPassword(username, password string) (*dotagiftx.Auth, error) {
	return s.findOne(dotagiftx.Auth{Username: username, Password: password})
}

func (s *authStorage) GetByRefreshToken(refreshToken string) (*dotagiftx.Auth, error) {
	row := &dotagiftx.Auth{}
	q := s.table().GetAllByIndex(authFieldRefreshToken, refreshToken)
	if err := s.db.one(q, row); err != nil {
		return nil, err
	}

	return row, nil
}

func (s *authStorage) Create(in *dotagiftx.Auth) error {
	t := now()
	in.CreatedAt = t
	in.UpdatedAt = t
	id, err := s.db.insert(s.table().Insert(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}
	in.ID = id

	return nil
}

func (s *authStorage) Update(in *dotagiftx.Auth) error {
	cur, err := s.Get(in.ID)
	if err != nil {
		return err
	}

	in.UpdatedAt = now()
	err = s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *authStorage) find(o dotagiftx.FindOpts) ([]dotagiftx.Auth, error) {
	var res []dotagiftx.Auth
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *authStorage) findOne(filter dotagiftx.Auth) (*dotagiftx.Auth, error) {
	o := dotagiftx.FindOpts{Filter: filter, Limit: 1}
	res, err := s.find(o)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, dotagiftx.AuthErrNotFound
	}

	return &res[0], nil
}

func (s *authStorage) table() r.Term {
	return r.Table(tableAuth)
}
