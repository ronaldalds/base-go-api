package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/i18n"
	"github.com/ronaldalds/base-go-api/internal/schemas"
	"github.com/ronaldalds/base-go-api/internal/utils"
	"github.com/ronaldalds/base-go-api/internal/validators"
)

func (con *Controller) ListUserHandler(c *fiber.Ctx) error {
	validate := validators.NewValidator()
	var params schemas.ListUsersRequest
	if err := c.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf(i18n.ERR_INVALID_QUERY_STRING, err.Error()),
		)
	}

	if validationErrors := validate.ValidateStruct(&params); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	users, err := con.Service.ListUsers()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	countUsers := uint(len(users))

	if err := utils.Pagination(params.Page, params.Limit, &users); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, i18n.ERR_PAGE_NOT_EXIST)
	}

	data := []schemas.UserResponse{}
	for _, user := range users {
		schemaUserResponse := &schemas.UserResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			Email:       user.Email,
			Active:      user.Active,
			IsSuperUser: user.IsSuperUser,
			RoleNames:   utils.ExtractNameRolesByUser(&user),
			Phone1:      user.Phone1,
			Phone2:      user.Phone2,
		}
		data = append(data, *schemaUserResponse)
	}

	res := &schemas.ListUsersResponse{
		Page:  params.Page,
		Limit: params.Limit,
		Data:  data,
		Total: countUsers,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (con *Controller) CreateUserHandler(c *fiber.Ctx) error {
	validate := validators.NewValidator()
	var body schemas.CreateUser
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(i18n.ERR_INVALID_BODY, err.Error()))
	}
	if validationErrors := validate.ValidateStruct(&body); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}
	if err := validate.ValidatePassword(body.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(i18n.ERR_CRYPTING_PASSWORD_FAILED, err.Error()))
	}
	body.Password = hashedPassword

	creator, err := utils.GetJwtHeaderPayload(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := con.Service.CreateUser(creator.Claims.Sub, body)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	res := &schemas.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Active:      user.Active,
		IsSuperUser: user.IsSuperUser,
		RoleNames:   utils.ExtractNameRolesByUser(user),
		Phone1:      user.Phone1,
		Phone2:      user.Phone2,
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (con *Controller) UpdateUserHandler(c *fiber.Ctx) error {
	validate := validators.NewValidator()
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, i18n.ERR_INVALID_ID_PARAMS)
	}
	var body schemas.UpdateUser
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(i18n.ERR_INVALID_BODY, err.Error()))
	}
	if validationErrors := validate.ValidateStruct(&body); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}
	editor, err := utils.GetJwtHeaderPayload(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	user, err := con.Service.UpdateUser(editor.Claims.Sub, uint(id), body)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	res := &schemas.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Active:      user.Active,
		IsSuperUser: user.IsSuperUser,
		RoleNames:   utils.ExtractNameRolesByUser(user),
		Phone1:      user.Phone1,
		Phone2:      user.Phone2,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
