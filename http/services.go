package http

import (
	"context"
	"io"
	"net/http"

	"github.com/kudarap/dotagiftx"
)

// authService provides access to auth service methods used by http handlers.
type authService interface {
	// SteamLogin redirects for authorization and process creation of auth.
	SteamLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*dotagiftx.Auth, error)

	// RevokeRefreshToken invalidates refresh token that will prevent on renewing
	// short-lived access token and will result user have to re-login.
	RevokeRefreshToken(ctx context.Context, refreshToken string) error

	// RefreshToken checks refresh token validity that allows to get new short-lived access token.
	RefreshToken(ctx context.Context, refreshToken string) (*dotagiftx.Auth, error)
}

// imageService provides access to image service methods used by http handlers.
type imageService interface {
	// Upload saves image details and actual file to local file system.
	Upload(context.Context, io.Reader) (fileID string, err error)

	// Image returns image details by id.
	Image(ctx context.Context, fileID string) (path string, err error)

	// Thumbnail downscales an image preserving its aspect ratio to the maximum dimensions.
	Thumbnail(ctx context.Context, fileID string, width, height uint) (path string, err error)
}

// itemService provides access to item service methods used by http handlers.
type itemService interface {
	// Items returns a list of items.
	Items(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Item, *dotagiftx.FindMetadata, error)

	// Item returns item details by id.
	Item(ctx context.Context, id string) (*dotagiftx.Item, error)

	// Create saves new item details.
	Create(context.Context, *dotagiftx.Item) error

	// Import creates new item from yaml format.
	Import(ctx context.Context, f io.Reader) (dotagiftx.ItemImportResult, error)

	// TopOrigins returns a list of top origin/treasure base on view count.
	TopOrigins(ctx context.Context) ([]string, error)

	// TopHeroes returns a list of top heroes base on view count.
	TopHeroes(ctx context.Context) ([]string, error)
}

// marketService provides access to market service methods used by http handlers.
type marketService interface {
	// Markets returns a list of markets.
	Markets(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Market, *dotagiftx.FindMetadata, error)

	// Market returns market details by id.
	Market(ctx context.Context, id string) (*dotagiftx.Market, error)

	// Create saves new market details.
	Create(context.Context, *dotagiftx.Market) error

	// Update saves market details changes.
	Update(context.Context, *dotagiftx.Market) error

	// Catalog returns a list of catalogs.
	Catalog(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Catalog, *dotagiftx.FindMetadata, error)

	// CatalogDetails returns catalog details by item id.
	CatalogDetails(ctx context.Context, id string, opts dotagiftx.FindOpts) (*dotagiftx.Catalog, error)

	// TrendingCatalog returns a top 10 trending catalogs.
	TrendingCatalog(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Catalog, *dotagiftx.FindMetadata, error)
}

// trackService provides access to track service methods used by http handlers.
type trackService interface {
	// CreateFromRequest saves new track from http request. Primarily used on client side.
	CreateFromRequest(ctx context.Context, r *http.Request) error

	// CreateSearchKeyword saves new keyword tracking data.
	CreateSearchKeyword(ctx context.Context, r *http.Request, keyword string) error
}

// statsService provides access to stats service methods used by http handlers.
type statsService interface {
	// CountMarketStatusV2 returns market status count base on given options.
	CountMarketStatusV2(ctx context.Context, opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error)

	// GraphMarketSales returns market sales graph base on given options.
	GraphMarketSales(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.MarketSalesGraph, error)

	// TopKeywords returns a list of top search keywords.
	TopKeywords(ctx context.Context) ([]dotagiftx.SearchKeywordScore, error)
}

// reportService provides access to report service methods used by http handlers.
type reportService interface {
	// Reports returns a list of reports.
	Reports(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Report, *dotagiftx.FindMetadata, error)

	// Report returns report details by id.
	Report(ctx context.Context, id string) (*dotagiftx.Report, error)

	// Create saves new report details.
	Create(context.Context, *dotagiftx.Report) error
}

// hammerService provides access to hammer service methods used by http handlers.
type hammerService interface {
	// Ban updates user status to ban and cancels all listings.
	Ban(context.Context, dotagiftx.HammerParams) (*dotagiftx.User, error)

	// Suspend updates user status to suspend and cancels all listings.
	Suspend(context.Context, dotagiftx.HammerParams) (*dotagiftx.User, error)

	// Lift update user status to "marked" and remove its ban or suspend a flag
	// and will restore items if requested.
	Lift(ctx context.Context, steamID string, restoreListings bool) error
}

// steamClient provides access to steam API methods used by http handlers.
type steamClient interface {
	// Player returns player summary base on steamID.
	Player(steamID string) (*dotagiftx.SteamPlayer, error)

	// ResolveVanityURL returns steam id from profile url.
	ResolveVanityURL(url string) (steamID string, err error)
}

// userService provides access to user service methods used by http handlers.
type userService interface {
	// FlaggedUsers returns a list of flagged/reported users.
	FlaggedUsers(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.User, error)

	// User returns user details by id.
	User(ctx context.Context, id string) (*dotagiftx.User, error)

	// UserFromContext returns user details from context.
	UserFromContext(context.Context) (*dotagiftx.User, error)

	// CreateSubscription creates a subscription for the current user.
	CreateSubscription(ctx context.Context, planID string) (subscriptionID string, err error)

	// ProcessSubscription validates and processes subscription features.
	ProcessSubscription(ctx context.Context, subscriptionID string) (*dotagiftx.User, error)

	// UpdateSubscriptionFromWebhook handles user subscription updates from http request.
	UpdateSubscriptionFromWebhook(ctx context.Context, r *http.Request) (*dotagiftx.User, error)

	// ProcessManualSubscription processes manual subscription.
	ProcessManualSubscription(ctx context.Context, form dotagiftx.ManualSubscriptionParam) (*dotagiftx.User, error)
}
