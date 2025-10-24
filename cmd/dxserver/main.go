package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/clickhouse"
	"github.com/kudarap/dotagiftx/config"
	"github.com/kudarap/dotagiftx/discord"
	"github.com/kudarap/dotagiftx/file"
	"github.com/kudarap/dotagiftx/http"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/paypal"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/tracing"
	"github.com/sirupsen/logrus"
)

const configPrefix = "DG"

var logger = logging.Default()

func main() {
	app := newApp()

	v := dotagiftx.NewVersion(false, tag, commit, built)
	logger.Println("version:", v.Tag)
	logger.Println("hash:", v.Commit)
	logger.Println("built:", v.Built)

	logger.Println("loading config...")
	if err := app.loadConfig(); err != nil {
		logger.Fatalln("could not load config:", err)
	}

	logger.Println("setting up...")
	if err := app.setup(); err != nil {
		logger.Fatalln("could not setup:", err)
	}

	logger.Println("running app...")
	if err := app.run(); err != nil {
		logger.Fatalln("could not run:", err)
	}
	logger.Println("stopped!")
}

type application struct {
	config config.Config
	server *http.Server
	logger *logrus.Logger

	closerFn func()
}

func (app *application) loadConfig() error {
	config.EnvPrefix = configPrefix
	if err := config.Load(&app.config); err != nil {
		return fmt.Errorf("load config: %s", err)
	}
	return nil
}

func (app *application) setup() error {
	// Logs setup.
	slogger := slog.Default()
	logger.Println("setting up persistent logs...")
	logSvc, err := logging.New(app.config.Log)
	if err != nil {
		return fmt.Errorf("could not set up logs: %s", err)
	}
	app.logger = logSvc

	// Database setup.
	logSvc.Println("setting up database...")
	redisClient, err := setupRedis(app.config.Redis)
	if err != nil {
		return err
	}
	clickHouseClient, err := setupClickHouse(app.config.ClickHouse)
	if err != nil {
		return err
	}
	rethinkClient, err := setupRethink(app.config.Rethink)
	if err != nil {
		return err
	}
	traceSpan := tracing.NewTracer(app.config.SpanEnabled, rethink.NewSpan(rethinkClient))
	rethinkClient.SetTracer(traceSpan)

	// Analytics stats capture.
	if err = setupChangeFeeds(rethinkClient, clickHouseClient); err != nil {
		return err
	}

	// External services setup.
	logSvc.Println("setting up external services...")
	steamClient, err := setupSteam(app.config.Steam, redisClient)
	if err != nil {
		return err
	}
	paypalClient, err := setupPaypal(app.config.Paypal)
	if err != nil {
		return err
	}
	discordClient := discord.New(app.config.DiscordWebhookURL)

	// Storage inits.
	logSvc.Println("setting up data stores...")
	userStg := rethink.NewUser(rethinkClient)
	authStg := rethink.NewAuth(rethinkClient)
	catalogStg := rethink.NewCatalog(rethinkClient, app.contextLog("storage_catalog"))
	itemStg := rethink.NewItem(rethinkClient)
	marketStg := rethink.NewMarket(rethinkClient)
	trackStg := rethink.NewTrack(rethinkClient)

	statsStg := rethink.NewStats(rethinkClient, app.contextLog("storage_stats"))
	reportStg := rethink.NewReport(rethinkClient)
	deliveryStg := rethink.NewDelivery(rethinkClient)
	inventoryStg := rethink.NewInventory(rethinkClient)

	// Service inits.
	logSvc.Println("setting up services...")
	fileMgr := setupFileManager(app.config)
	userSvc := dotagiftx.NewUserService(userStg, fileMgr, paypalClient)
	authSvc := dotagiftx.NewAuthService(app.config.SigKey, steamClient, authStg, userSvc)
	imageSvc := dotagiftx.NewImageService(fileMgr)
	itemSvc := dotagiftx.NewItemService(app.config.AllowedImageSources, itemStg, fileMgr)
	inventorySvc := dotagiftx.NewInventoryService(inventoryStg, marketStg, catalogStg)
	deliverySvc := dotagiftx.NewDeliveryService(deliveryStg, marketStg)
	marketSvc := dotagiftx.NewMarketService(
		marketStg,
		userStg,
		itemStg,
		trackStg,
		catalogStg,
		statsStg,
		deliverySvc,
		inventorySvc,
		steamClient,
		rethink.NewQueue(rethinkClient),
		app.contextLog("service_market"),
	)
	trackSvc := dotagiftx.NewTrackService(trackStg, itemStg)
	reportSvc := dotagiftx.NewReportService(reportStg, discordClient)
	statsSvc := dotagiftx.NewStatsService(statsStg, trackStg)
	hammerSvc := dotagiftx.NewHammerService(userStg, marketStg)
	phantasmSvc := phantasm.NewService(app.config.Phantasm, redisClient, slogger)

	// Server setup.
	logSvc.Println("setting up http server...")
	srv := http.NewServer(
		app.config.SigKey,
		app.config.DivineKey,
		userSvc,
		authSvc,
		imageSvc,
		itemSvc,
		marketSvc,
		trackSvc,
		statsSvc,
		reportSvc,
		hammerSvc,
		steamClient,
		phantasmSvc,
		traceSpan,
		redisClient,
		initVer(app.config),
		logSvc,
	)
	srv.Addr = app.config.Addr
	app.server = srv

	app.closerFn = func() {
		logSvc.Println("closing and stopping app...")
		if err = redisClient.Close(); err != nil {
			logSvc.Fatal("could not close redis client", err)
		}
		if err = rethinkClient.Close(); err != nil {
			logSvc.Fatal("could not close rethink client", err)
		}
		if err = clickHouseClient.Close(); err != nil {
			logSvc.Fatal("could not close clickhouse client", err)
		}
	}

	return nil
}

func (app *application) run() error {
	defer app.closerFn()
	return app.server.Run()
}

func (app *application) contextLog(name string) logging.Logger {
	return logging.WithPrefix(app.logger, name)
}

func newApp() *application {
	a := &application{}
	a.closerFn = func() {}
	return a
}

func setupSteam(cfg steam.Config, rc *redis.Client) (*steam.Client, error) {
	c, err := steam.New(cfg, rc)
	if err != nil {
		return nil, fmt.Errorf("could not setup steam client: %s", err)
	}

	return c, nil
}

func setupPaypal(cfg paypal.Config) (*paypal.Client, error) {
	c, err := paypal.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not setup paypal client: %s", err)
	}

	return c, nil
}

func setupFileManager(cfg config.Config) *file.Local {
	c := cfg.Upload
	return file.New(c.Path, c.Size, c.Types)
}

func setupRethink(cfg rethink.Config) (c *rethink.Client, err error) {
	c = &rethink.Client{}
	fn := func() error {
		c, err = rethink.New(cfg)
		if err != nil {
			return fmt.Errorf("could not setup rethink client: %s", err)
		}

		return nil
	}

	err = connRetry("rethink", fn)
	return
}

func setupRedis(cfg redis.Config) (c *redis.Client, err error) {
	c = &redis.Client{}
	fn := func() error {
		c, err = redis.New(cfg)
		if err != nil {
			return fmt.Errorf("could not setup redis client: %s", err)
		}

		return nil
	}

	err = connRetry("redis", fn)
	return
}

func setupClickHouse(cfg clickhouse.Config) (c *clickhouse.Client, err error) {
	c = &clickhouse.Client{}
	fn := func() error {
		c, err = clickhouse.New(cfg)
		if err != nil {
			return fmt.Errorf("could not setup clickhouse client: %s", err)
		}

		return nil
	}

	err = connRetry("clickhouse", fn)
	return
}

func setupChangeFeeds(rethinkClient *rethink.Client, clickhouseClient *clickhouse.Client) error {
	ctx := context.Background()
	err := rethinkClient.ListenChangeFeed("track", func(prev, next []byte) error {
		var v dotagiftx.Track
		if err := json.Unmarshal(next, &v); err != nil {
			return err
		}
		return clickhouseClient.CaptureTrackStats(ctx, v)
	})
	if err != nil {
		return err
	}

	err = rethinkClient.ListenChangeFeed("market", func(prev, next []byte) error {
		fmt.Println("market change feed", prev, next)

		deleted := prev != nil && next == nil
		if deleted {
			var v dotagiftx.Market
			if err := json.Unmarshal(prev, &v); err != nil {
				return err
			}
			return clickhouseClient.DeleteMarketStats(ctx, v.ID)
		}

		if next == nil {
			return nil
		}
		var v dotagiftx.Market
		if err := json.Unmarshal(next, &v); err != nil {
			return err
		}
		return clickhouseClient.CaptureMarketStats(ctx, v)
	})
	return err
}

func connRetry(name string, fn func() error) error {
	const delay = time.Second * 5

	// Catches a panic to retry.
	defer func() {
		if err := recover(); err != nil {
			logger.Printf("[%s] conn error: %s. retrying in %s...", name, err, delay)
			time.Sleep(delay)
			_ = connRetry(name, fn)
		}
	}()

	// Causes panic to retry.
	if err := fn(); err != nil {
		panic(err)
	}

	return nil
}

// version details used by ldflags.
var tag, commit, built string

func initVer(cfg config.Config) *dotagiftx.Version {
	v := dotagiftx.NewVersion(cfg.Prod, tag, commit, built)
	return v
}
