package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/config"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/worker"
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
	worker *worker.Worker
	logger *logrus.Logger

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
	logger.Println("setting up persistent logs...")
	logSvc, err := logging.New(app.config.Log)
	if err != nil {
		return fmt.Errorf("could not set up logs: %s", err)
	}
	app.logger = logSvc

	// Database setup.
	logSvc.Println("setting up database...")
	rethinkClient, err := setupRethink(app.config.Rethink)
	if err != nil {
		return err
	}
	defer rethinkClient.Close()

	migrate(rethinkClient.Session())
	os.Exit(0)

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

func (app *application) contextLog(name string) logging.Logger {
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
