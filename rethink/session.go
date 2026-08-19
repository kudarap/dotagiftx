package rethink

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableAuthSession         = "auth_session"
	sessionFieldRefreshToken = "refresh_token"
	sessionFieldExpiresAt    = "expires_at"
)

// NewSession creates a new instance of session data store.
func NewSession(c *Client) *SessionRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableAuthSession); err != nil {
		log.Fatalf("could not create %s table: %s", tableAuthSession, err)
	}

	if err := c.autoIndex(ctx, tableAuthSession, dotagiftx.AuthSession{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableAuthSession, err)
	}

	return &SessionRepository{c}
}

type SessionRepository struct {
	db *Client
}

func (s *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*dotagiftx.AuthSession, error) {
	row := &dotagiftx.AuthSession{}
	q := s.table().GetAllByIndex(sessionFieldRefreshToken, refreshToken).Limit(1)
	if err := s.db.one(ctx, q, row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.AuthErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *SessionRepository) Create(ctx context.Context, in *dotagiftx.AuthSession) error {
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

func (s *SessionRepository) Update(ctx context.Context, in *dotagiftx.AuthSession) error {
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

func (s *SessionRepository) Delete(ctx context.Context, id string) error {
	return s.db.delete(ctx, s.table().Get(id).Delete())
}

func (s *SessionRepository) Get(ctx context.Context, id string) (*dotagiftx.AuthSession, error) {
	row := &dotagiftx.AuthSession{}
	if err := s.db.one(ctx, s.table().Get(id), row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.AuthErrNotFound
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *SessionRepository) PurgeExpiredSessions(ctx context.Context, cutOff time.Time) error {
	q := s.table().Filter(r.Row.Field(sessionFieldExpiresAt).Lt(cutOff)).Delete()
	if err := s.db.exec(ctx, q); err != nil {
		return fmt.Errorf("deleting expired session: %w", err)
	}

	return nil
}

func (s *SessionRepository) table() r.Term {
	return r.Table(tableAuthSession)
}
