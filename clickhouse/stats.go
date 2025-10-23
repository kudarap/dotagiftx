package clickhouse

import (
	"context"
	"fmt"

	"github.com/kudarap/dotagiftx"
)

func (c *Client) CaptureTrack(ctx context.Context, track dotagiftx.Track) error {
	panic("implement me")
}

func (c *Client) CaptureMarket(ctx context.Context, market dotagiftx.Market) error {
	const qry = `INSERT INTO market (
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

	err := c.conn.AsyncInsert(
		ctx,
		qry,
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
