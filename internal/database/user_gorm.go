package database

import (
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ronaldalds/base-go-api/internal/models"
	"github.com/ronaldalds/base-go-api/internal/settings"
	"github.com/ronaldalds/base-go-api/internal/utils"
	"github.com/ronaldalds/base-go-api/internal/validators"
	"gorm.io/gorm"
)

func (gs *GormStore) CreateAdmin() error {
	validate := validators.NewValidator()
	var user models.User
	err := gs.DB.Where("username = ?", settings.Env.SuperUsername).First(&user).Error
	if err == nil {
		return fmt.Errorf("admin already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check admin existence: %s", err.Error())
	}
	if err := validate.ValidatePassword(settings.Env.SuperPass); err != nil {
		return fmt.Errorf("failed to validate password")
	}
	hashPassword, err := utils.HashPassword(settings.Env.SuperPass)
	if err != nil {
		return fmt.Errorf("failed to create admin: %s", err.Error())
	}
	admin := &models.User{
		FirstName:   settings.Env.SuperName,
		LastName:    "Admin",
		Username:    settings.Env.SuperUsername,
		Email:       settings.Env.SuperEmail,
		Password:    hashPassword,
		Active:      true,
		IsSuperUser: true,
		Phone1:      settings.Env.SuperPhone,
	}
	if err := gs.DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create user: %s", err.Error())
	}
	return fmt.Errorf("admin created successfully")
}

func (gs *GormStore) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := gs.DB.
		Preload("Roles.Permissions").
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no record found for id: %d", id)
		}
		return nil, fmt.Errorf("failed to query database: %w", result.Error)
	}
	return &user, nil
}

func (gs *GormStore) GetUserByUsernameOrEmail(text string) (*models.User, error) {
	var user models.User
	result := gs.DB.
		Preload("Roles.Permissions").
		Where("username = ? OR email = ?", text, text).
		First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no record found for email: %s", text)
		}
		return nil, fmt.Errorf("failed to query database: %w", result.Error)
	}
	return &user, nil
}

func (gs *GormStore) GetUsers() ([]models.User, error) {
	var users []models.User
	result := gs.DB.
		Preload("Roles.Permissions").
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query database: %w", result.Error)
	}
	return users, nil
}

func (gs *GormStore) CheckIfUserExistsByUsernameOrEmail(email, username string) error {
	var user models.User
	if err := gs.DB.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return fmt.Errorf("no record found for email or username")
	}
	return nil
}
