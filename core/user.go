package core

import (
	"context"
	"time"
)

// User error types.
const (
	UserErrNotFound Errors = iota + 1100
	UserErrRequiredID
	UserErrRequiredFields
	UserErrProfileImageDL
	UserErrSteamSync
)

// sets error text definition.
func init() {
	appErrorText[UserErrNotFound] = "user not found"
	appErrorText[UserErrRequiredID] = "user id is required"
	appErrorText[UserErrRequiredFields] = "user fields are required"
	appErrorText[UserErrProfileImageDL] = "user profile image could not download"
	appErrorText[UserErrSteamSync] = "user profile steam sync error"
}

type (
	// User represents user information.
	User struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		SteamID   string     `json:"steam_id"   db:"steam_id,omitempty"    valid:"required"`
		Name      string     `json:"name"       db:"name,omitempty"        valid:"required"`
		URL       string     `json:"url"        db:"url,omitempty"         valid:"required"`
		Avatar    string     `json:"avatar"     db:"avatar,omitempty"      valid:"required"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
	}

	// UserService provides access to user service.
	UserService interface {
		// Users returns a list of users.
		Users(opts FindOpts) ([]User, error)

		// User returns user details by id.
		User(id string) (*User, error)

		// Create saves new user and download profile image to local file.
		Create(*User) error

		// UserFromContext returns user details from context.
		UserFromContext(context.Context) (*User, error)

		// Update saves user changes.
		Update(context.Context, *User) error

		// SteamSync saves updated steam info.
		SteamSync(sp *SteamPlayer) (*User, error)
	}

	// UserStorage defines operation for user records.
	UserStorage interface {
		// Find returns a list of users from data store.
		Find(opts FindOpts) ([]User, error)

		// Get returns user details by id from data store.
		Get(id string) (*User, error)

		// Create persists a new user to data store.
		Create(*User) error

		// Update persists user changes to data store.
		Update(*User) error
	}
)

// CheckCreate validates field on creating new user.
func (u User) CheckCreate() error {
	return validator.Struct(u)
}

// CheckUpdate validates field on update user.
func (u User) CheckUpdate() error {
	if u.ID == "" {
		return UserErrRequiredID
	}

	return nil
}
