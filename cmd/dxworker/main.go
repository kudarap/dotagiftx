package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kudarap/dotagiftx/config"
	"github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/tracing"
	"github.com/kudarap/dotagiftx/verify"
	"github.com/kudarap/dotagiftx/worker"
	"github.com/kudarap/dotagiftx/worker/jobs"
)

const configPrefix = "DG"

var logger = logging.Default()

func main() {
	app := newApp()

	v := dotagiftx.NewVersion(false, tag, commit, built)
	logger.Info("version", "tag", v.Tag)
	logger.Info("hash", "commit", v.Commit)
	logger.Info("built", "built", v.Built)

	logger.Info("loading config...")
	if err := app.loadConfig(); err != nil {
		logger.Error("could not load config", "error", err)
		os.Exit(1)
	}

	logger.Info("setting up...")
	if err := app.setup(); err != nil {
		logger.Error("could not setup", "error", err)
		os.Exit(1)
	}

	logger.Info("running app...")
	if err := app.run(); err != nil {
		logger.Error("could not run", "error", err)
		os.Exit(1)
	}
	logger.Info("stopped!")
}

type application struct {
	config config.Config
	server *nethttp.Server
	worker *worker.Worker
	logger *slog.Logger

	closerFn func()
}

func (app *application) loadConfig() error {
	config.EnvPrefix = configPrefix
	if err := config.Load(&app.config); err != nil {
		return fmt.Errorf("could not load config: %s", err)
	}

	return nil
}

func (app *application) setup() error {
	// Logs setup.
	logger.Info("setting up persistent logs...")
	logSvc, err := logging.New(app.config.Log)
	if err != nil {
		return fmt.Errorf("could not set up logs: %s", err)
	}
	app.logger = logSvc

	// Database setup.
	logSvc.Info("setting up database...")
	redisClient, err := setupRedis(app.config.Redis)
	if err != nil {
		return err
	}
	rethinkClient, err := setupRethink(app.config.Rethink)
	if err != nil {
		return err
	}
	traceSpan := tracing.NewTracer(app.config.SpanEnabled, rethink.NewSpan(rethinkClient))
	rethinkClient.SetTracer(traceSpan)

	// External services setup.
	logSvc.Info("setting up external services...")

	// Storage inits.
	logSvc.Info("setting up data stores...")
	catalogStg := rethink.NewCatalog(rethinkClient, app.contextLog("storage_catalog"))
	marketStg := rethink.NewMarket(rethinkClient)
	deliveryStg := rethink.NewDelivery(rethinkClient)
	inventoryStg := rethink.NewInventory(rethinkClient)
	userStg := rethink.NewUser(rethinkClient)
	sessionStg := rethink.NewSession(rethinkClient)
	queue := rethink.NewQueue(rethinkClient)

	// Service inits.
	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slogger := slog.New(th)
	logSvc.Info("setting up services...")
	inventorySvc := dotagiftx.NewInventoryService(inventoryStg, marketStg, catalogStg)
	deliverySvc := dotagiftx.NewDeliveryService(deliveryStg, marketStg)
	phantasmSvc := phantasm.NewService(app.config.Phantasm, redisClient, slogger)
	verifySources := []verify.AssetSource{phantasmSvc.InventoryAssetWithProvider}
	// TODO: Use proper level of fallbacks. For experimental purposes only.
	if len(app.config.Phantasm.BackupAddrs) != 0 {
		logSvc.Info("EXPERIMENTAL: phantasm backup source enabled", "backup_addrs", app.config.Phantasm.BackupAddrs)
		phantasmSvcExp := phantasm.NewService(phantasm.Config{
			Addrs:                app.config.Phantasm.BackupAddrs,
			WebhookURL:           app.config.Phantasm.WebhookURL,
			Secret:               app.config.Phantasm.Secret,
			Path:                 app.config.Phantasm.Path,
			MaxFetchRetryAttempt: 1,
		}, redisClient, slogger)
		verifySources = append(verifySources, phantasmSvcExp.InventoryAssetWithProvider)
	}
	verifySources = append(verifySources, steaminvorg.InventoryAssetWithProvider)
	assetSource := verify.NewSource(verifySources...)

	// Setup application worker
	tp := worker.NewTaskProcessor(time.Second, queue, inventorySvc, deliverySvc, assetSource, phantasmSvc, slogger)
	app.worker = worker.New(tp)
	app.worker.SetLogger(app.contextLog("worker"))
	app.worker.AddJob(jobs.NewRecheckInventory(
		inventorySvc,
		marketStg,
		assetSource,
		logging.WithPrefix(logger, "job_recheck_inventory"),
	))
	app.worker.AddJob(jobs.NewVerifyInventory(
		inventorySvc,
		marketStg,
		assetSource,
		logging.WithPrefix(logger, "job_verify_inventory"),
	))
	app.worker.AddJob(jobs.NewVerifyDelivery(
		deliverySvc,
		marketStg,
		assetSource,
		logging.WithPrefix(logger, "job_verify_delivery"),
	))
	app.worker.AddJob(jobs.NewGiftWrappedUpdate(
		deliverySvc,
		deliveryStg,
		marketStg,
		assetSource,
		logging.WithPrefix(logger, "job_giftwrapped_update"),
	))
	app.worker.AddJob(jobs.NewRevalidateDelivery(
		deliverySvc,
		marketStg,
		assetSource,
		logging.WithPrefix(logger, "job_revalidate_delivery"),
	))
	app.worker.AddJob(jobs.NewExpiringSubscription(
		userStg,
		redisClient,
		logging.WithPrefix(logger, "job_expiring_subscription"),
	))
	app.worker.AddJob(jobs.NewExpiringMarket(
		marketStg,
		catalogStg,
		redisClient,
		logging.WithPrefix(logger, "job_expiring_market"),
	))
	app.worker.AddJob(jobs.NewSweepMarket(marketStg, logging.WithPrefix(logger, "job_sweep_market")))
	app.worker.AddJob(jobs.NewSweepPhantasmCache(phantasmSvc, logging.WithPrefix(logger, "job_sweep_phantasm")))
	app.worker.AddJob(jobs.NewSweepAuthSess(sessionStg, logging.WithPrefix(logger, "job_sweep_auth_sess")))

	// Server setup.
	logSvc.Info("setting up http server...")
	app.server = setupServer(app.config.Addr, phantasmSvc)

	app.closerFn = func() {
		logSvc.Info("closing and stopping app...")
		// Force server shutdown after shutdownTimeout and this was added because of SSE's opened connection.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err = app.server.Shutdown(ctx); err != nil {
			logSvc.Error("could not shutdown http server", "error", err)
		}
		if err = app.worker.Stop(); err != nil {
			logSvc.Error("could not stop worker", "error", err)
			os.Exit(1)
		}
		if err = redisClient.Close(); err != nil {
			logSvc.Error("could not close redis client", "error", err)
			os.Exit(1)
		}
		if err = rethinkClient.Close(); err != nil {
			logSvc.Error("could not close rethink client", "error", err)
			os.Exit(1)
		}
	}

	return nil
}

func (app *application) run() error {
	defer app.closerFn()

	// Handle quit on SIGINT (CTRL-C).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Handle error on server start.
	svlog := app.contextLog("server")
	go func() {
		svlog.Info("starting server on", "addr", app.config.Addr)
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			svlog.Error("could not listen and serve", "addr", app.config.Addr, "error", err)
			os.Exit(1)
		}
	}()

	// delay worker start to give leeway on phantasm webhook to be online
	time.Sleep(5 * time.Second)
	go app.worker.Start()

	<-quit
	return nil
}

func (app *application) contextLog(name string) *slog.Logger {
	return logging.WithPrefix(app.logger, name)
}

func newApp() *application {
	a := &application{}
	a.closerFn = func() {}
	return a
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

func setupServer(addr string, svc *phantasm.Service) *nethttp.Server {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Post("/webhook/phantasm/{steam_id}", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		id := chi.URLParam(r, "steam_id")
		secret := r.Header.Get(phantasm.WebhookAuthHeader)
		if err := svc.SaveInventory(r.Context(), id, secret, r.Body); err != nil {
			w.WriteHeader(nethttp.StatusOK)
			_, _ = fmt.Fprintf(w, "error: %s", err)
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		_, _ = fmt.Fprintf(w, "ok")
	})

	const readWriteTimeout = time.Second * 15
	return &nethttp.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: readWriteTimeout,
		ReadTimeout:  readWriteTimeout,
	}
}

func connRetry(name string, fn func() error) error {
	const delay = time.Second * 5

	// Catches a panic to retry.
	defer func() {
		if err := recover(); err != nil {
			logger.Info("conn error, retrying", "name", name, "error", err, "delay", delay)
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
