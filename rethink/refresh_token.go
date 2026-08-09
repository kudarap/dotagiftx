package rethink

import (
	"context"
	"errors"
	"log"

	"dario.cat/mergo"
	"github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableRefreshToken          = "refresh_token"
	refreshTokenFieldTokenHash = "token_hash"
	refreshTokenFieldFamilyID  = "family_id"
)

// NewRefreshToken creates a new instance of refresh token data store.
func NewRefreshToken(c *Client) *RefreshTokenRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableRefreshToken); err != nil {
		log.Fatalf("could not create %s table: %s", tableRefreshToken, err)
	}

	if err := c.autoIndex(ctx, tableRefreshToken, dotagiftx.RefreshToken{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableRefreshToken, err)
	}

	return &RefreshTokenRepository{c}
}

type RefreshTokenRepository struct {
	db *Client
}

func (s *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*dotagiftx.RefreshToken, error) {
	row := &dotagiftx.RefreshToken{}
	q := s.table().GetAllByIndex(refreshTokenFieldTokenHash, tokenHash)
	if err := s.db.one(ctx, q, row); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return nil, dotagiftx.AuthErrRefreshToken
		}

		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return row, nil
}

func (s *RefreshTokenRepository) Create(ctx context.Context, in *dotagiftx.RefreshToken) error {
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

func (s *RefreshTokenRepository) Update(ctx context.Context, in *dotagiftx.RefreshToken) error {
	cur := &dotagiftx.RefreshToken{}
	if err := s.db.one(ctx, s.table().Get(in.ID), cur); err != nil {
		if errors.Is(err, r.ErrEmptyResult) {
			return dotagiftx.AuthErrRefreshToken
		}

		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	in.UpdatedAt = now()
	err := s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	if err := mergo.Merge(in, cur); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageMergeErr, err)
	}

	return nil
}

func (s *RefreshTokenRepository) RevokeFamily(ctx context.Context, familyID string) error {
	q := s.table().GetAllByIndex(refreshTokenFieldFamilyID, familyID).Update(map[string]any{
		"revoked":    true,
		"updated_at": now(),
	})
	if err := s.db.update(ctx, q); err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return nil
}

func (s *RefreshTokenRepository) table() r.Term {
	return r.Table(tableRefreshToken)
}
