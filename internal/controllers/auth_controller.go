package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/handlers"
	"github.com/ronaldalds/base-go-api/internal/i18n"
	"github.com/ronaldalds/base-go-api/internal/schemas"
	"github.com/ronaldalds/base-go-api/internal/settings"
	"github.com/ronaldalds/base-go-api/internal/utils"
	"github.com/ronaldalds/base-go-api/internal/validators"
)

func (con *Controller) LoginHandler(c *fiber.Ctx) error {
	errors := handlers.NewError()
	validate := validators.NewValidator()
	var body schemas.LoginRequest
	// validate request body
	if err := c.BodyParser(&body); err != nil {
		errors.AddDetailErr("bodyParser", fmt.Sprintf(i18n.ERR_INVALID_BODY, err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	if validationErrors := validate.ValidateStruct(&body); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	// find username or email in database
	user, err := con.Service.Login(body)
	if err != nil {
		errors.AddDetailErr("Login", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(errors)
	}
	// check password
	if !utils.CheckPasswordHash(body.Password, user.Password) {
		errors.AddDetailErr("Login", i18n.ERR_LOGIN)
		return c.Status(fiber.StatusUnauthorized).JSON(errors)
	}
	// generate tokens
	token, err := utils.GenerateToken(user, settings.Env.JwtExpireAcess)
	if err != nil {
		errors.AddDetailErr("GenerateAccessToken", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(errors)
	}
	// save token in redis
	if err := con.Service.SetToken(user.ID, token); err != nil {
		errors.AddDetailErr("SetTokenRedis", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(errors)
	}
	// send response
	res := &schemas.LoginResponse{
		Token: token,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
