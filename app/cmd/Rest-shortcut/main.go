package main

import (
	"Rest-shortcut/internal/config"
	"Rest-shortcut/internal/http-server/handlers/login"
	"Rest-shortcut/internal/http-server/handlers/redirect"
	"Rest-shortcut/internal/http-server/handlers/register"
	"Rest-shortcut/internal/http-server/handlers/save"
	"Rest-shortcut/internal/http-server/middlewares/authMiddleware"
	"Rest-shortcut/lib/api"
	sl "Rest-shortcut/lib/logger"
	"Rest-shortcut/storage/postrges"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	conf := config.NewConfig()

	logger := setupLogger(conf.Environment)

	storage, err := postrges.NewStorage(conf.Storage)
	if err != nil {
		logger.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("Connect to database")
	auth := setupAuth(conf.AuthConfig)
	router := setupRouter(logger, storage, auth)
	logger.Info("Setup handlers")
	_ = router
	logger.Info("Start server", slog.String("address", conf.HttpConfig.Address+conf.HttpConfig.Port))
	server := &http.Server{
		Addr:    conf.HttpConfig.Address + conf.HttpConfig.Port,
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server", sl.Err(err))
	}
	logger.Info("Server shutdown")
}

func setupAuth(authConfig config.AuthConfig) *api.Auth {
	auth := api.NewAuth(authConfig.AccessTtl, authConfig.RefreshTtl, authConfig.SignedKey)
	return auth
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("Initialize logger")
	return logger
}

func setupRouter(logger *slog.Logger, storage *postrges.Storage, auth *api.Auth) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/register", register.New(logger, storage))
	router.Get("/login", login.New(logger, storage, auth))
	router.Post("/refresh", login.New(logger, storage, auth))
	router.Route("/user", func(r chi.Router) {
		r.Use(authMiddleware.New(logger, auth))
		r.Post("/", save.New(logger, storage))
		r.Get("/{text}", redirect.New(logger, storage))
	})
	return router
}
