package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ronaldalds/base-go-api/internal/controllers"
	"github.com/ronaldalds/base-go-api/internal/middlewares"
)

type Router struct {
	App        *fiber.App
	Controller *controllers.Controller
	Middleware *middlewares.Middleware
}

func NewRouter(app *fiber.App) *Router {
	return &Router{
		App:        app,
		Controller: controllers.NewController(),
		Middleware: middlewares.NewMiddleware(app),
	}
}

func (r *Router) RegisterFiberRoutes() {
	r.Middleware.CorsMiddleware()
	r.Middleware.SecurityMiddleware()
	apiV2 := r.App.Group("/api/v2")
	apiV2.Get("/health", r.Controller.HealthHandler)

	// Group authentication
	authGroup := apiV2.Group("/auth", r.Middleware.Limited(10))
	r.Auth(authGroup)

	// Group Users
	usersGroup := apiV2.Group("/users")
	r.User(usersGroup)
	r.Role(usersGroup)
	r.Permission(usersGroup)
}
