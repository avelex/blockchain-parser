package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/avelex/blockchain-parser/config"
	"github.com/avelex/blockchain-parser/internal/api"
	"github.com/avelex/blockchain-parser/internal/ethclient"
	"github.com/avelex/blockchain-parser/internal/parser"
	"github.com/avelex/blockchain-parser/internal/repository/memory"
)

var configFile = flag.String("config", "config.yaml", "Config file path")

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	context.AfterFunc(ctx, func() {
		slog.Info("Interrupt signal received, Stopping blockchain parser...")
	})

	defer cancel()

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("Loaded config", "port", cfg.Port, "rpc", cfg.RPC, "blocks_interval", cfg.BlocksInterval, "start_block", cfg.StartBlock)

	repo := memory.New()
	client := ethclient.New(cfg.RPC)
	parser := parser.New(cfg, client, repo)
	handler := api.NewHandler(parser)

	mux := http.NewServeMux()
	handler.Register(mux)

	go func() {
		slog.Info("Starting Blockchain Parser")

		if err := parser.Start(ctx); err != nil {
			slog.Error("Failed to start parser", "error", err)
		}
	}()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		slog.Info("Starting HTTP server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start http server", "error", err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Warn("Failed to shutdown http server", "error", err)
	}

	slog.Info("Http server stopped")

	slog.Info("Blockchain parser stopped")
}
