package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/handlers"
)

// New cria uma nova instância do FiberServer, inicializando o Fiber e o banco.
func New() *fiber.App {
	// Cria o servidor Fiber encapsulado na estrutura
	server := fiber.New(fiber.Config{
		ServerHeader: "R.A.L.D.S",
		ErrorHandler: handlers.ErrorHandler,
	})
	return server
}
