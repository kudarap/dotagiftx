package verified

import "github.com/kudarap/dotagiftx/steam"

type VerifyStatus string

const (
	VerifyStatusError   VerifyStatus = "error"
	VerifyStatusPrivate VerifyStatus = "private"
	VerifyStatusNoHit   VerifyStatus = "no-hit"
	VerifyStatusItem    VerifyStatus = "item"
	VerifyStatusSeller  VerifyStatus = "seller"
)

// AssetSource represents inventory asset source provider.
type AssetSource func(steamID string) ([]steam.Asset, error)
