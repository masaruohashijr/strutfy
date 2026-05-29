package app

import (
	"context"
	appauth "github.com/masaru/strutfy-api/internal/auth"
	"github.com/masaru/strutfy-api/internal/config"
	"github.com/masaru/strutfy-api/internal/database"
	apphttp "github.com/masaru/strutfy-api/internal/http"
	"github.com/masaru/strutfy-api/internal/http/handlers"
	"github.com/masaru/strutfy-api/internal/organization"
	"github.com/masaru/strutfy-api/internal/project"
	"github.com/masaru/strutfy-api/internal/service"
	"github.com/masaru/strutfy-api/internal/user"
	"github.com/masaru/strutfy-api/internal/workflow"
	"log"
	"net/http"
	"time"
)

type App struct {
	cfg    config.Config
	server *http.Server
	engine *workflow.Engine
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	engine := workflow.NewEngine(100)

	tokenService := service.NewTokenService(cfg.JWTSecret, cfg.JWTIssuer)
	userRepo := user.NewRepository(db)
	organizationRepo := organization.NewRepository(db)
	authService := appauth.NewService(
		userRepo,
		organizationRepo,
		tokenService,
	)
	authHandler := handlers.NewAuthHandler(authService)

	projectRepo := project.NewRepository(db)
	projectService := project.NewService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService, engine)

	router := apphttp.NewRouter(apphttp.RouterDeps{
		AuthHandler:    authHandler,
		ProjectHandler: projectHandler,
		Engine:         engine,
		TokenService:   tokenService,
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		cfg:    cfg,
		server: server,
		engine: engine,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.engine.Start(ctx, 4)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = a.server.Shutdown(shutdownCtx)
	}()

	log.Printf("Strutfy API running on port %s", a.cfg.HTTPPort)

	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
