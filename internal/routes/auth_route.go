package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) Auth(router fiber.Router) {
	router.Post("/login", r.Controller.LoginHandler)
	// router.Post("/logout", r.Controller.LogoutHandler)
}
