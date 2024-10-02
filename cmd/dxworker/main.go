package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kudarap/dotagiftx/gokit/envconf"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/gokit/version"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/service"
	"github.com/kudarap/dotagiftx/tracing"
	"github.com/kudarap/dotagiftx/worker"
	"github.com/kudarap/dotagiftx/worker/jobs"
	"github.com/sirupsen/logrus"
)

const configPrefix = "DG"

var logger = log.Default()

func main() {
	app := newApp()

	v := version.New(false, tag, commit, built)
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
	config  Config
	worker  *worker.Worker
	logger  *logrus.Logger
	version *version.Version

	closerFn func()
}

func (app *application) loadConfig() error {
	envconf.EnvPrefix = configPrefix
	if err := envconf.Load(&app.config); err != nil {
		return fmt.Errorf("could not load config: %s", err)
	}

	return nil
}

func (app *application) setup() error {
	// Logs setup.
	logger.Println("setting up persistent logs...")
	logSvc, err := log.New(app.config.Log)
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
	rethinkClient, err := setupRethink(app.config.Rethink)
	if err != nil {
		return err
	}
	traceSpan := tracing.NewTracer(app.config.SpanEnabled, rethink.NewSpan(rethinkClient))
	rethinkClient.SetTracer(traceSpan)

	// External services setup.
	logSvc.Println("setting up external services...")

	// Storage inits.
	logSvc.Println("setting up data stores...")
	catalogStg := rethink.NewCatalog(rethinkClient, app.contextLog("storage_catalog"))
	marketStg := rethink.NewMarket(rethinkClient)
	deliveryStg := rethink.NewDelivery(rethinkClient)
	inventoryStg := rethink.NewInventory(rethinkClient)
	userStg := rethink.NewUser(rethinkClient)
	queue := rethink.NewQueue(rethinkClient)

	// Service inits.
	logSvc.Println("setting up services...")
	deliverySvc := service.NewDelivery(deliveryStg, marketStg)
	inventorySvc := service.NewInventory(inventoryStg, marketStg, catalogStg)

	// Setup application worker
	tp := worker.NewTaskProcessor(time.Second, queue, inventorySvc, deliverySvc)
	app.worker = worker.New(tp)
	app.worker.SetLogger(app.contextLog("worker"))
	app.worker.AddJob(jobs.NewRecheckInventory(
		inventorySvc, marketStg, log.WithPrefix(logger, "job_recheck_inventory"),
	))
	app.worker.AddJob(jobs.NewVerifyInventory(
		inventorySvc,
		marketStg,
		log.WithPrefix(logger, "job_verify_inventory"),
	))
	app.worker.AddJob(jobs.NewVerifyDelivery(
		deliverySvc,
		marketStg,
		log.WithPrefix(logger, "job_verify_delivery"),
	))
	app.worker.AddJob(jobs.NewGiftWrappedUpdate(
		deliverySvc,
		deliveryStg,
		marketStg,
		log.WithPrefix(logger, "job_giftwrapped_update"),
	))
	app.worker.AddJob(jobs.NewRevalidateDelivery(
		deliverySvc,
		marketStg,
		log.WithPrefix(logger, "job_revalidate_delivery"),
	))
	app.worker.AddJob(jobs.NewExpiringSubscription(
		userStg,
		redisClient,
		log.WithPrefix(logger, "job_expiring_subscription"),
	))
	app.worker.AddJob(jobs.NewExpiringMarket(
		marketStg,
		catalogStg,
		redisClient,
		log.WithPrefix(logger, "job_expiring_market"),
	))
	app.worker.AddJob(jobs.NewSweepMarket(
		marketStg, log.WithPrefix(logger, "job_sweep_market"),
	))

	app.closerFn = func() {
		logSvc.Println("closing and stopping app...")
		if err = app.worker.Stop(); err != nil {
			logSvc.Fatal("could not stop worker", err)
		}
		if err = redisClient.Close(); err != nil {
			logSvc.Fatal("could not close redis client", err)
		}
		if err = rethinkClient.Close(); err != nil {
			logSvc.Fatal("could not close rethink client", err)
		}
	}

	return nil
}

func (app *application) run() error {
	defer app.closerFn()

	// Handle quit on SIGINT (CTRL-C).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go app.worker.Start()

	<-quit
	return nil
}

func (app *application) contextLog(name string) log.Logger {
	return log.WithPrefix(app.logger, name)
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

func connRetry(name string, fn func() error) error {
	const delay = time.Second * 5

	// Catches a panic to retry.
	defer func() {
		if err := recover(); err != nil {
			logger.Printf("[%s] conn error: %s. retrying in %s...", name, err, delay)
			time.Sleep(delay)
			err = connRetry(name, fn)
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
