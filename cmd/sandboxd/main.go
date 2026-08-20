package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/application"
	"github.com/example/wasm-sandbox-runtime/internal/config"
	"github.com/example/wasm-sandbox-runtime/internal/repository/memory"
	"github.com/example/wasm-sandbox-runtime/internal/runtime/mock"
	httptransport "github.com/example/wasm-sandbox-runtime/internal/transport/http"
	"github.com/example/wasm-sandbox-runtime/internal/worker"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := memory.NewStore()
	rt := mock.NewRuntime()
	svc := application.NewService(store, rt, cfg, log)
	queue := worker.NewScheduler(svc, cfg.Workers)
	queue.Start()
	server := &http.Server{Addr: cfg.Address, Handler: httptransport.New(svc, log)}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server", "error", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	queue.Stop()
	_ = server.Shutdown(ctx)
}
