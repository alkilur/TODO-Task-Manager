package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"yet-another-todo-list/internal/http-server"
	"yet-another-todo-list/internal/lib/slwrap"
	"yet-another-todo-list/internal/storage/sqlite"

	"yet-another-todo-list/internal/config"
	"yet-another-todo-list/internal/lib/slogpretty"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{},
	}
	log := slog.New(opts.NewPrettyHandler(os.Stdout))
	log.Info(fmt.Sprintf("config: %s", cfg))

	// init storage
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Error("failed to init storage", slwrap.Wrap(err))
		os.Exit(1)
	}
	_ = db

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	controller := http_server.New()
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.Get("/api/nextdate", controller.GetNextDate(log))

	// run server
	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err)
	}
	log.Error("server stopped")
}
