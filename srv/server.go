package srv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/middleware"
)

func NewServer(
	logger *log.Logger,
	config *config.Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, config)
	var handler http.Handler = mux
	handler = middleware.PerformanceLogger(handler, logger)
	handler = middleware.AddState(handler)
	handler = middleware.CorrelationId(handler)
	return handler
}

func Run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := log.New(os.Stdout, "[logicng-service] ", log.LstdFlags)
	server := NewServer(logger, cfg)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: server,
	}
	go func() {
		logger.Printf("listening on %s with sync timeout of %s\n", httpServer.Addr, cfg.SyncComputationTimout)
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
