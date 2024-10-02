package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/http/jwt"
	gokitMw "github.com/kudarap/dotagiftx/gokit/http/middleware"
	"github.com/kudarap/dotagiftx/gokit/version"
	"github.com/kudarap/dotagiftx/tracing"
	"github.com/sirupsen/logrus"
)

const (
	defaultAddr      = ":8000"
	shutdownTimeout  = 10 * time.Second
	readWriteTimeout = 15 * time.Second
)

// NewServer returns new http server.
func NewServer(
	sigKey string,
	divineKey string,
	us dgx.UserService,
	au dgx.AuthService,
	is dgx.ImageService,
	its dgx.ItemService,
	ms dgx.MarketService,
	ts dgx.TrackService,
	ss dgx.StatsService,
	rs dgx.ReportService,
	hs dgx.HammerService,
	sc dgx.SteamClient,
	t *tracing.Tracer,
	c dgx.Cache,
	v *version.Version,
	l *logrus.Logger,
) *Server {
	jwt.SigKey = sigKey
	return &Server{
		divineKey: divineKey,
		userSvc:   us,
		authSvc:   au,
		imageSvc:  is,
		itemSvc:   its,
		marketSvc: ms,
		trackSvc:  ts,
		statsSvc:  ss,
		reportSvc: rs,
		hammerSvc: hs,
		steam:     sc,
		tracing:   t,
		cache:     c,
		logger:    l,
		version:   v,
	}
}

// Server represents http Server.
type Server struct {
	// Server settings.
	Addr    string
	handler http.Handler
	// Service resources.
	userSvc   dgx.UserService
	authSvc   dgx.AuthService
	imageSvc  dgx.ImageService
	itemSvc   dgx.ItemService
	marketSvc dgx.MarketService
	trackSvc  dgx.TrackService
	statsSvc  dgx.StatsService
	reportSvc dgx.ReportService
	hammerSvc dgx.HammerService
	steam     dgx.SteamClient

	tracing *tracing.Tracer
	cache   dgx.Cache
	logger  *logrus.Logger
	version *version.Version

	// divineKey is a special access key for importing and creating items and
	// managing manual subscriptions. This key is used as temporary admin access key.
	divineKey string
}

func (s *Server) setup() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(s.tracing.Middleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(NewStructuredLogger(s.logger))
	r.Use(gokitMw.CORS)
	r.Use(middleware.Recoverer)

	// Set routes handler.
	s.publicRouter(r)
	s.privateRouter(r)

	r.NotFound(handle404())
	r.MethodNotAllowed(handle405())

	// Set default address.
	if s.Addr == "" {
		s.Addr = defaultAddr
	}

	s.handler = r
}

func (s *Server) Run() error {
	s.setup()

	// Setup http router.
	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      s.handler,
		WriteTimeout: readWriteTimeout,
		ReadTimeout:  readWriteTimeout,
	}

	// Handle error on server start.
	errCh := make(chan error, 1)
	go func() {
		s.logger.Infoln("server running on", s.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Handle quit on SIGINT (CTRL-C).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Force server shutdown after shutdownTimeout and this was added because of SSE's opened connection.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	select {
	case err := <-errCh:
		return err
	case <-quit:
		s.logger.Infoln("server shutting down...")
		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Error("server shutdown error", err)
		}
		s.logger.Infoln("server stopped!")
		return nil
	}
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger})
}

func isValidDivineKey(r *http.Request, divineKey string) error {
	if r.URL.Query().Get("key") == divineKey {
		return nil
	}
	return fmt.Errorf("divine key does not exist or invalid")
}
