package controllers

import (
	"github.com/ronaldalds/base-go-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	Service *services.Service
}

func NewController() *Controller {
	return &Controller{
		Service: services.NewService(),
	}
}

func (con *Controller) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(con.Service.Health())
}
