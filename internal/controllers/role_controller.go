package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/handlers"
	"github.com/ronaldalds/base-go-api/internal/i18n"
	"github.com/ronaldalds/base-go-api/internal/schemas"
	"github.com/ronaldalds/base-go-api/internal/utils"
	"github.com/ronaldalds/base-go-api/internal/validators"
)

func (con *Controller) ListRoleHandler(c *fiber.Ctx) error {
	validate := validators.NewValidator()
	var params schemas.ListRoleRequest
	if err := c.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(i18n.ERR_INVALID_QUERY_STRING, err.Error()))
	}

	if validationErrors := validate.ValidateStruct(&params); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	roles, err := con.Service.ListRoles()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	countRoles := uint(len(roles))

	if err := utils.Pagination(params.Page, params.Limit, &roles); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, i18n.ERR_PAGE_NOT_EXIST)
	}

	data := []schemas.RoleResponse{}
	for _, role := range roles {
		schemaRoleResponse := &schemas.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
		for _, permission := range role.Permissions {
			schemaRoleResponse.Permissions = append(schemaRoleResponse.Permissions, schemas.PermissionResponse{
				ID:          permission.ID,
				Code:        permission.Code,
				Name:        permission.Name,
				Description: permission.Description,
			})
		}
		data = append(data, *schemaRoleResponse)
	}

	res := &schemas.ListRoleResponse{
		Page:  params.Page,
		Limit: params.Limit,
		Data:  data,
		Total: countRoles,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (con *Controller) CreateRoleHandler(c *fiber.Ctx) error {
	errors := handlers.NewError()
	validate := validators.NewValidator()
	var body schemas.CreateRoleRequest
	if err := c.BodyParser(&body); err != nil {
		errors.AddDetailErr("bodyParser", fmt.Sprintf(i18n.ERR_INVALID_BODY, err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	if validationErrors := validate.ValidateStruct(&body); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}
	role, err := con.Service.CreateRole(body)
	if err != nil {
		errors.AddDetailErr("Register", err.Error())
		return c.Status(fiber.StatusCreated).JSON(errors)
	}
	// Extrair os codes das permissões
	var permissionCodes []schemas.PermissionResponse
	for _, permission := range role.Permissions {
		schemaPermission := &schemas.PermissionResponse{
			ID:          permission.ID,
			Code:        permission.Code,
			Name:        permission.Name,
			Description: permission.Description,
		}
		permissionCodes = append(permissionCodes, *schemaPermission)
	}

	// Preparar a resposta
	res := &schemas.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissionCodes, // Adicionar apenas os codes
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
