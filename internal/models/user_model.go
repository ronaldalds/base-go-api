package models

import "time"

// User representa o modelo de usuário no sistema.
type User struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	FirstName   string    `gorm:"size:50;not null" validate:"required,min=1,max=50"` // Nome obrigatório
	LastName    string    `gorm:"size:50" validate:"omitempty,max=50"`               // Sobrenome opcional
	Username    string    `gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
	Email       string    `gorm:"uniqueIndex;not null" validate:"required,email"`
	Password    string    `gorm:"not null" validate:"required"`                        // Hash da senha
	Active      bool      `gorm:"default:true"`                                        // Se o usuário está ativo
	IsSuperUser bool      `gorm:"default:false"`                                       // Se o usuário é superuser
	Roles       []Role    `gorm:"many2many:user_roles"`                                // Relacionamento com grupos
	Phone1      string    `gorm:"type:varchar(20);not null" validate:"required,e164"`  // Telefone obrigatório (formato E.164)
	Phone2      string    `gorm:"type:varchar(20);nullable" validate:"omitempty,e164"` // Telefone opcional (formato E.164)
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement"`
	Name        string       `gorm:"uniqueIndex;size:100;not null" validate:"required,min=3,max=100"`
	Description string       `gorm:"size:255"`                    // Descrição opcional do grupo
	Permissions []Permission `gorm:"many2many:roles_permissions"` // Relacionamento com permissões
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
}

type Permission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" validate:"required"`
	Name        string    `gorm:"uniqueIndex;size:100;not null" validate:"required,min=3,max=100"`
	Description string    `gorm:"size:255"` // Descrição opcional da permission
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
