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
	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/phantasm"
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
	us dotagiftx.UserService,
	au dotagiftx.AuthService,
	is dotagiftx.ImageService,
	its dotagiftx.ItemService,
	ms dotagiftx.MarketService,
	ts dotagiftx.TrackService,
	ss dotagiftx.StatsService,
	rs dotagiftx.ReportService,
	hs dotagiftx.HammerService,
	sc dotagiftx.SteamClient,
	ps *phantasm.Service,
	t *tracing.Tracer,
	c Cache,
	v *dotagiftx.Version,
	l *logrus.Logger,
) *Server {
	SigKey = sigKey
	return &Server{
		divineKey:   divineKey,
		userSvc:     us,
		authSvc:     au,
		imageSvc:    is,
		itemSvc:     its,
		marketSvc:   ms,
		trackSvc:    ts,
		statsSvc:    ss,
		reportSvc:   rs,
		hammerSvc:   hs,
		steam:       sc,
		phantasmSvc: ps,
		tracing:     t,
		cache:       c,
		logger:      l,
		version:     v,
	}
}

// Server represents http Server.
type Server struct {
	// Server settings.
	Addr    string
	handler http.Handler
	// Service resources.
	userSvc   dotagiftx.UserService
	authSvc   dotagiftx.AuthService
	imageSvc  dotagiftx.ImageService
	itemSvc   dotagiftx.ItemService
	marketSvc dotagiftx.MarketService
	trackSvc  dotagiftx.TrackService
	statsSvc  dotagiftx.StatsService
	reportSvc dotagiftx.ReportService
	hammerSvc dotagiftx.HammerService
	steam     dotagiftx.SteamClient

	phantasmSvc *phantasm.Service

	tracing *tracing.Tracer
	cache   Cache
	logger  *logrus.Logger
	version *dotagiftx.Version

	// divineKey is a special access key for importing and creating items and
	// managing manual subscriptions. This key is used as temporary admin access key.
	divineKey string
}

func (s *Server) setup() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(s.tracing.Middleware)
	r.Use(middleware.RequestID)
	r.Use(vercelRequestID)
	r.Use(middleware.RealIP)
	r.Use(NewStructuredLogger(s.logger))
	r.Use(cors)
	r.Use(requestIDWriter)
	r.Use(middleware.Recoverer)

	// Set routes handler.
	s.publicRouter(r)
	s.privateRouter(r)

	r.NotFound(handle404())
	r.MethodNotAllowed(handle405())

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
