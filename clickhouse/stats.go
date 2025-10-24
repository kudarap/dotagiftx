package clickhouse

import (
	"context"
	"fmt"

	"github.com/kudarap/dotagiftx"
)

func (c *Client) CaptureTrackStats(ctx context.Context, track dotagiftx.Track) error {
	const query = `INSERT INTO track (
		id,
		type,
		item_id,
		user_id,
		keyword,
		client_ip,
		user_agent,
		referer,
		cookies,
		sess_user_id,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	err := c.db.AsyncInsert(
		ctx,
		query,
		false,
		track.ID,
		track.Type,
		track.ItemID,
		track.UserID,
		track.Keyword,
		track.ClientIP,
		track.UserAgent,
		track.Referer,
		track.Cookies,
		track.SessUserID,
		track.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("async insert track: %w", err)
	}

	return nil
}

func (c *Client) CaptureMarketStats(ctx context.Context, market dotagiftx.Market) error {
	const query = `INSERT INTO market (
		id,
		user_id,
		item_id,
		type,
		status,
		price,
		currency,
		partner_steam_id,
		inventory_status,
		delivery_status,
		resell,
		seller_steam_id,
		created_at,
		updated_at,
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	err := c.db.AsyncInsert(
		ctx,
		query,
		false,
		market.ID,
		market.UserID,
		market.ItemID,
		uint(market.Type),
		uint(market.Status),
		market.Price,
		market.Currency,
		market.PartnerSteamID,
		uint(market.InventoryStatus),
		uint(market.DeliveryStatus),
		market.Resell,
		market.SellerSteamID,
		market.CreatedAt.Unix(),
		market.UpdatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("async insert market: %w", err)
	}

	return nil
}

func (c *Client) DeleteMarketStats(ctx context.Context, id string) error {
	// delete is too slow
	return nil

	err := c.db.Exec(ctx, `DELETE FROM market WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete market: %w", err)
	}
	return nil
}

func (c *Client) CountMarketStatus(ctx context.Context, opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	panic("implement me")
}

func (c *Client) TopKeywords(ctx context.Context) ([]dotagiftx.SearchKeywordScore, error) {
	panic("implement me")
}

func (c *Client) TrendingCatalog(
	ctx context.Context,
	opts dotagiftx.FindOpts,
) ([]dotagiftx.Catalog, *dotagiftx.FindMetadata, error) {
	panic("implement me")
}
