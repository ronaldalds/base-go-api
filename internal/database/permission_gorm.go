package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ronaldalds/base-go-api/internal/models"
	"gorm.io/gorm"
)

func findMissingIDsByPermissions(ids []uint, permissions []models.Permission) []uint {
	// Criar um mapa dos IDs encontrados
	foundIDs := make(map[uint]struct{})
	for _, p := range permissions {
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

func (gs *GormStore) GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	result := gs.DB.
		Find(&permissions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query database: %w", result.Error)
	}
	return permissions, nil
}

func (gs *GormStore) GetPermissionByIds(ids []uint) ([]models.Permission, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("no permission IDs provided")
	}

	var permissions []models.Permission
	// Buscar as permissões pelos IDs fornecidos
	if err := gs.DB.Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch permissions: %s", err.Error())
	}

	// Verificar se todas as permissões foram encontradas
	if len(permissions) != len(ids) {
		missingIDs := findMissingIDsByPermissions(ids, permissions)
		return nil, fmt.Errorf("permissions not found for IDs: %v", missingIDs)
	}

	return permissions, nil
}

func (gs *GormStore) CheckIfPermissionExistsByCodeOrName(code, name string) error {
	var permission models.Permission
	result := gs.DB.Where("code = ? OR name = ?", code, name).First(&permission)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("no record found for code or name")
		}
		return fmt.Errorf("failed to query database: %w", result.Error)
	}
	return nil
}
