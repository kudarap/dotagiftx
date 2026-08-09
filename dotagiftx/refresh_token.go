package dotagiftx

import (
	"context"
	"time"
)

// DefaultRefreshTokenLifetime defines the time to live of a refresh token.
const DefaultRefreshTokenLifetime = 30 * 24 * time.Hour

// RefreshToken represents a single refresh token session. A family of tokens
// share the same FamilyID to enable reuse detection and family wide revocation.
type RefreshToken struct {
	ID        string     `json:"id"         db:"id,omitempty"`
	AuthID    string     `json:"auth_id"    db:"auth_id,indexed,omitempty"`
	UserID    string     `json:"user_id"    db:"user_id,indexed,omitempty"`
	TokenHash string     `json:"-"          db:"token_hash,indexed,omitempty"`
	FamilyID  string     `json:"-"          db:"family_id,indexed,omitempty"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at,omitempty"`
	Revoked   bool       `json:"-"          db:"revoked,omitempty"`
	CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
}

// refreshTokenRepository defines operations for refresh token records.
type refreshTokenRepository interface {
	// GetByTokenHash returns a refresh token by its hashed value.
	GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)

	// Create persists a new refresh token.
	Create(context.Context, *RefreshToken) error

	// Update persists refresh token changes.
	Update(context.Context, *RefreshToken) error

	// RevokeFamily invalidates all tokens belonging to a token family.
	RevokeFamily(ctx context.Context, familyID string) error
}
