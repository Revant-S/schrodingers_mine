package main

import (
	"Schrodingers_Mine/api"
	"Schrodingers_Mine/config"
	"Schrodingers_Mine/solver"
	"Schrodingers_Mine/solver/deduce"
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// precedence: flags > real environment > .env > defaults
	if err := config.LoadDotEnv(".env"); err != nil {
		log.Fatalf("load .env: %v", err)
	}
	defBudget, err := config.Duration("BOT_MOVE_BUDGET", 4*time.Second)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	defShutdown, err := config.Duration("BOT_SHUTDOWN_TIMEOUT", 3*time.Second)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	defInsecure, err := config.Bool("BOT_INSECURE_SKIP_VERIFY", true)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	addr := flag.String("addr", config.String("BOT_ADDR", "localhost:8080"), "listen address (keep it on localhost)")
	budget := flag.Duration("budget", defBudget, "thinking time per move; keep below the board's move timeout (default 5s)")
	name := flag.String("name", config.String("BOT_NAME", "schrodingers_mine"), "bot name sent in READY")
	shutdownTimeout := flag.Duration("shutdown-timeout", defShutdown, "grace period for open connections on shutdown")
	insecure := flag.Bool("insecure-skip-verify", defInsecure, "accept WebSocket connections from any origin")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	solve := &deduce.Deduce{}
	mux.Handle("/", api.Handler(api.Config{
		Name:               *name,
		MoveBudget:         *budget,
		InsecureSkipVerify: *insecure,
		NewSolver:          func() solver.Solver { return solve }, // swap in your solver here
	}))

	srv := &http.Server{
		Addr:        *addr,
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	go func() {
		log.Printf("bot listening on ws://%s", *addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
