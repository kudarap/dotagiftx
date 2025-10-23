package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kudarap/dotagiftx"
	"github.com/schollz/progressbar/v3"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var (
	minBatchRows     int64 = 8
	ratePerIteration       = 1000
)

func init() {
	flag.Int64Var(&minBatchRows, "minbatchrows", 8, "MinBatchRows")
	flag.IntVar(&ratePerIteration, "rateperiteration", 1000, "RatePerIteration")
	flag.Parse()
}

func migrate(db *r.Session) {
	ctx := context.Background()

	cursor, err := r.Table("market_production").
		OrderBy(r.OrderByOpts{Index: "created_at"}).
		Run(db, r.RunOpts{
			MinBatchRows: minBatchRows,
		})
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	chdb := clickhouseInit()

	var total int
	bench := time.Now()
	progressBar := progressbar.Default(120000)

	var batchInsertCtr = 0
	for {
		var market dotagiftx.Market
		if cursor.Next(&market) {
			total++

			// Clickhouse writes
			if ratePerIteration != 0 {
				batchInsertCtr++
				if batchInsertCtr == ratePerIteration {
					time.Sleep(1 * time.Second)
					batchInsertCtr = 0
				}
			}

			err = chdb.AsyncInsert(
				ctx,
				`INSERT INTO market (
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
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
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
				log.Println("error append clickhouse:", total, err)
			}

			progressBar.Add(1)
			continue
		}

		break
	}

	fmt.Println("total:", total, "took:", time.Since(bench))
}

func clickhouseInit() driver.Conn {
	conn, err := connect()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT name, toString(uuid) as uuid_str FROM system.tables LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name, uuid string
		if err := rows.Scan(&name, &uuid); err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, uuid: %s", name, uuid)
	}

	return conn
}

func connect() (driver.Conn, error) {
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{},
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to clickhouse: %s", err)
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, fmt.Errorf("could not ping clickhouse: %s", err)
	}
	return conn, nil
}
