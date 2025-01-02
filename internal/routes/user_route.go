package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) User(router fiber.Router) {
	router.Get(
		"/",
		r.Middleware.JWTProtected("view_user"),
		r.Controller.ListUserHandler,
	)
	router.Post(
		"/",
		r.Middleware.JWTProtected("create_user"),
		r.Controller.CreateUserHandler,
	)
	router.Put(
		"/:id",
		r.Middleware.JWTProtected(),
		r.Controller.UpdateUserHandler,
	)
}

func (r *Router) Role(router fiber.Router) {
	router.Get(
		"/roles",
		r.Middleware.JWTProtected(),
		r.Controller.ListRoleHandler,
	)
	router.Post(
		"/roles",
		r.Middleware.JWTProtected("create_role"),
		r.Controller.CreateRoleHandler,
	)
}

func (r *Router) Permission(router fiber.Router) {
	router.Get(
		"/permissions",
		r.Middleware.JWTProtected(),
		r.Controller.ListPermissiontHandler,
	)
}
