package middlewares

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/utils"
)

// This middleware needs a Redis client to check tokens
func (m *Middleware) JWTProtected(permissions ...string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token, err := utils.GetJwtHeaderPayload(c)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Check Redis for active session
		session, err := m.RedisStore.GetKey(fmt.Sprintf("%d", token.Claims.Sub))
		if err != nil {
			log.Println("Redis error:", err) // Log the Redis error
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		if session != token.Token {
			log.Println("Token does not match active session") // Log the error
			return fiber.NewError(fiber.StatusUnauthorized, "token does not match active session")
		}

		// Check permissions
		user, err := m.GormStore.GetUserByID(token.Claims.Sub)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		if !user.Active {
			return fiber.NewError(fiber.StatusUnauthorized, "user is inactive")
		}
		if user.IsSuperUser {
			return c.Next()
		}
		if len(permissions) == 0 {
			return c.Next()
		}

		// Gather user permissions into a slice
		userPermissions := utils.ExtractCodePermissionsByUser(user)

		// Check if any required permission exists in user's permissions
		for _, requiredPermission := range permissions {
			if utils.Contains(userPermissions, requiredPermission) {
				log.Println("Permission validated, proceeding to next handler")
				return c.Next()
			}
		}

		// If no errors, log success and continue to the next handler
		log.Println("JWT validated and session matched, proceeding to next handler")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
}
