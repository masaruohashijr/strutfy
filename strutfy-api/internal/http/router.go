package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/masaru/strutfy-api/internal/http/handlers"
	"github.com/masaru/strutfy-api/internal/http/middleware"
	"github.com/masaru/strutfy-api/internal/service"
	"github.com/masaru/strutfy-api/internal/workflow"
)

type RouterDeps struct {
    AuthHandler    *handlers.AuthHandler
	ProjectHandler *handlers.ProjectHandler
	Engine         *workflow.Engine
	TokenService   *service.TokenService
}

func NewRouter(deps RouterDeps) http.Handler {
	router := httprouter.New()

	auth := func(next httprouter.Handle) httprouter.Handle {
		return middleware.Auth(deps.TokenService, next)
	}

	// Públicas
	router.GET("/health", health)

	// Auth
    router.POST("/v1/auth/register", deps.AuthHandler.Register)
    router.POST("/v1/auth/login", deps.AuthHandler.Login)
    router.POST("/v1/auth/refresh", notImplemented)

	// Projetos
	router.POST(
		"/v1/projects",
		auth(deps.ProjectHandler.Create),
	)

	router.DELETE(
		"/v1/projects/:id",
		auth(deps.ProjectHandler.Delete),
	)

	return router
}

func health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func notImplemented(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	_, _ = w.Write([]byte(`{"error":"not implemented"}`))
}