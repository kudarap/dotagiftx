package rethink

import (
	"log"

	"github.com/kudarap/dotagiftx/core"
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

	if err := c.autoIndex(tableAuth, core.Auth{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableAuth, err)
	}

	return &authStorage{c}
}

type authStorage struct {
	db *Client
}

func (s *authStorage) Get(id string) (*core.Auth, error) {
	row := &core.Auth{}
	if err := s.db.one(s.table().Get(id), row); err != nil {
		if err == r.ErrEmptyResult {
			return nil, core.AuthErrNotFound
		}

		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *authStorage) GetByUsername(username string) (*core.Auth, error) {
	row := &core.Auth{}
	q := s.table().GetAllByIndex(authFieldUsername, username)
	if err := s.db.one(q, row); err != nil {
		return nil, err
	}

	return row, nil
}

func (s *authStorage) GetByUsernameAndPassword(username, password string) (*core.Auth, error) {
	return s.findOne(core.Auth{Username: username, Password: password})
}

func (s *authStorage) GetByRefreshToken(refreshToken string) (*core.Auth, error) {
	row := &core.Auth{}
	q := s.table().GetAllByIndex(authFieldRefreshToken, refreshToken)
	if err := s.db.one(q, row); err != nil {
		return nil, err
	}

	return row, nil
}

func (s *authStorage) Create(in *core.Auth) error {
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

func (s *authStorage) Update(in *core.Auth) error {
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

func (s *authStorage) find(o core.FindOpts) ([]core.Auth, error) {
	var res []core.Auth
	q := newFindOptsQuery(s.table(), o)
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}

	return res, nil
}

func (s *authStorage) findOne(filter core.Auth) (*core.Auth, error) {
	o := core.FindOpts{Filter: filter, Limit: 1}
	res, err := s.find(o)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, core.AuthErrNotFound
	}

	return &res[0], nil
}

func (s *authStorage) table() r.Term {
	return r.Table(tableAuth)
}
