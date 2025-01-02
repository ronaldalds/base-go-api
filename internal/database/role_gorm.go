package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ronaldalds/base-go-api/internal/models"
	"gorm.io/gorm"
)

func findMissingIDsByRoles(ids []uint, roles []models.Role) []uint {
	// Criar um mapa dos IDs encontrados
	foundIDs := make(map[uint]struct{})
	for _, p := range roles {
		foundIDs[p.ID] = struct{}{}
	}

	// Identificar os IDs ausentes
	var missingIDs []uint
	for _, id := range ids {
		if _, exists := foundIDs[id]; !exists {
			missingIDs = append(missingIDs, id)
		}
	}

	return missingIDs
}

func (gs *GormStore) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	result := gs.DB.
		Preload("Permissions").
		Find(&roles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query database: %w", result.Error)
	}
	return roles, nil
}

func (gs *GormStore) GetRoleByIds(ids []uint) ([]models.Role, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("no role IDs provided")
	}

	var roles []models.Role
	// Buscar as permissões pelos IDs fornecidos
	if err := gs.DB.Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %s", err.Error())
	}

	// Verificar se todas as permissões foram encontradas
	if len(roles) != len(ids) {
		missingIDs := findMissingIDsByRoles(ids, roles)
		return nil, fmt.Errorf("roles not found for IDs: %v", missingIDs)
	}

	return roles, nil
}

func (gs *GormStore) CheckIfRoleExistsByIds(id []uint) error {
	var rolesCount int64
	if err := gs.DB.Model(&models.Role{}).
		Where("id IN ?", id).
		Count(&rolesCount).Error; err != nil {
		return fmt.Errorf("failed to validate roles: %s", err.Error())
	}
	if rolesCount != int64(len(id)) {
		return fmt.Errorf("some roles are invalid or do not exist")
	}
	return nil
}

func (gs *GormStore) CheckIfRoleExistsByName(name string) error {
	var role models.Role
	result := gs.DB.Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("no record found for name")
		}
		return fmt.Errorf("failed to query database: %w", result.Error)
	}
	return nil
}
