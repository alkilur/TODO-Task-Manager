package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Handle("/*", http.FileServer(http.Dir("../web")))

	// run server
	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err)
	}
	log.Error("server stopped")
}
