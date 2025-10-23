package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const initDatabase = `
	CREATE TABLE IF NOT EXISTS track
	(
		id UUID,
		type Enum8('' = 0, 'v' = 1, 's' = 2, 'p' = 3),
		item_id Nullable(String),
		user_id Nullable(String),
		keyword Nullable(String),
		client_ip Nullable(IPv4),
		user_agent Nullable(String),
		referer Nullable(String),
		cookies Array(Nullable(String)),
		sess_user_id Nullable(String),
		created_at DateTime
	)
	ENGINE = MergeTree
	PARTITION BY toYYYYMM(created_at)
	ORDER BY created_at;

	CREATE TABLE IF NOT EXISTS market
	(
		id UUID,
		user_id String,
		item_id String,
		type UInt16,
		status UInt16,
		inventory_status UInt16,
		delivery_status UInt16,
		price Float32,
		currency String,
		partner_steam_id String,
		resell Bool,
		seller_steam_id String,
		created_at DateTime,
		updated_at DateTime
	)
	ENGINE = ReplacingMergeTree
	PARTITION BY toYYYYMM(updated_at)
	ORDER BY id;`

func autoMigration(ctx context.Context, db driver.Conn) error {
	_, err := db.Query(ctx, initDatabase)
	return err
}
