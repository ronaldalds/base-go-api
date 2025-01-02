package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/i18n"
	"github.com/ronaldalds/base-go-api/internal/schemas"
	"github.com/ronaldalds/base-go-api/internal/utils"
	"github.com/ronaldalds/base-go-api/internal/validators"
)

func (con *Controller) ListPermissiontHandler(c *fiber.Ctx) error {
	validate := validators.NewValidator()
	var params schemas.ListPermissionRequest
	if err := c.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(i18n.ERR_INVALID_QUERY_STRING, err.Error()))
	}

	if validationErrors := validate.ValidateStruct(&params); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	permissions, err := con.Service.ListPermissions()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	countPermissions := uint(len(permissions))

	if err := utils.Pagination(params.Page, params.Limit, &permissions); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, i18n.ERR_PAGE_NOT_EXIST)
	}

	data := []schemas.PermissionResponse{}
	for _, permission := range permissions {
		schemaPermissionResponse := &schemas.PermissionResponse{
			ID:          permission.ID,
			Code:        permission.Code,
			Name:        permission.Name,
			Description: permission.Description,
		}
		data = append(data, *schemaPermissionResponse)
	}

	res := &schemas.ListPermissionResponse{
		Page:  params.Page,
		Limit: params.Limit,
		Data:  data,
		Total: countPermissions,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
