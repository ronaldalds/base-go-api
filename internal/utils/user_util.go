package utils

import "github.com/ronaldalds/base-go-api/internal/models"

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsAll(listX, listY []models.Role) bool {
	// Criar um mapa para os itens de X
	itemMap := make(map[uint]bool)
	for _, item := range listX {
		itemMap[item.ID] = true
	}

	// Verificar se todos os itens de Y estão no mapa de X
	for _, item := range listY {
		if !itemMap[item.ID] {
			return false // Item de Y não está em X
		}
	}

	return true // Todos os itens de Y estão em X
}

func ExtractNameRolesByUser(user *models.User) []string {
	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames
}

func ExtractCodePermissionsByUser(user *models.User) []string {
	var codePermissions []string
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			codePermissions = append(codePermissions, permission.Code)
		}
	}
	return codePermissions
}
