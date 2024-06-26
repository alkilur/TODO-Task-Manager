package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"yet-another-todo-list/internal/config"
	"yet-another-todo-list/internal/http-server/handlers"
	"yet-another-todo-list/internal/lib/slogpretty"
	"yet-another-todo-list/internal/lib/slwrap"
	"yet-another-todo-list/internal/storage/sqlite"

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

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.Get("/api/nextdate", handlers.GetNextDate(log))
	router.Post("/api/task", handlers.AddTask(log, db))
	router.Get("/api/tasks", handlers.GetTasks(log, db))
	router.Put("/api/task", handlers.UpdateTask(log, db))

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
