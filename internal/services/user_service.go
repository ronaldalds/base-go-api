package services

import (
	"fmt"

	"github.com/ronaldalds/base-go-api/internal/models"
	"github.com/ronaldalds/base-go-api/internal/schemas"
	"github.com/ronaldalds/base-go-api/internal/utils"
)

func (s *Service) CreateUser(creatorID uint, req schemas.CreateUser) (*models.User, error) {
	// Verificar se o username ou email já existe
	err := s.GormStore.CheckIfUserExistsByUsernameOrEmail(req.Email, req.Username)
	if err == nil {
		return nil, fmt.Errorf("user with username '%s' or email '%s' already exists", req.Username, req.Email)
	}

	// Buscar as roles pelo ID
	roles, err := s.GormStore.GetRoleByIds(req.Roles)
	if err != nil {
		return nil, fmt.Errorf("role with ids '%v' does not exist", req.Roles)
	}

	// Buscar o criador do usuário
	creator, err := s.GormStore.GetUserByID(creatorID)
	if err != nil {
		return nil, fmt.Errorf("user with id '%v' does not exist", creatorID)
	}

	// Validar se o criador possui as roles necessárias ou é superusuário
	if !creator.IsSuperUser && !utils.ContainsAll(creator.Roles, roles) {
		return nil, fmt.Errorf("failed to create user: creator does not have all required roles")
	}

	// Criar o usuário (apenas em memória)
	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Active:    req.Active,
		IsSuperUser: func() bool {
			if req.IsSuperUser {
				if creator.IsSuperUser {
					return true
				}
				panic(fmt.Errorf("only superusers can create other superusers"))
			}
			return false
		}(),
		Phone1: req.Phone1,
		Phone2: req.Phone2,
	}

	// Persistir o usuário no banco de dados
	if err := s.GormStore.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %s", err.Error())
	}

	// Associar as roles ao usuário
	if err := s.GormStore.DB.Model(&user).Association("Roles").Replace(roles); err != nil {
		return nil, fmt.Errorf("failed to set roles for user: %v", err)
	}

	// Retornar o usuário criado
	return &user, nil
}

func (s *Service) UpdateUser(editorID uint, id uint, req schemas.UpdateUser) (*models.User, error) {
	user, err := s.GormStore.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user modified with id '%v' does not exists", id)
	}
	editor, err := s.GormStore.GetUserByID(editorID)
	if err != nil {
		return nil, fmt.Errorf("user editor with id '%v' does not exist", editorID)
	}

	// Verificar permissões para alterar Roles
	canUpdateRoles := false
	if editor.IsSuperUser {
		canUpdateRoles = true
	} else {
		// Buscar permissões associadas às roles do editor
		permissions := utils.ExtractCodePermissionsByUser(editor)

		// Verificar se o editor possui a permissão `update_user`
		if utils.Contains(permissions, "update_user") {
			canUpdateRoles = true
		}
	}

	// Atualizar as Roles somente se permitido
	if len(req.Roles) > 0 {
		if !canUpdateRoles {
			return nil, fmt.Errorf("editor does not have permission to update user roles")
		}
		// Buscar as roles especificadas na atualização
		roles, err := s.GormStore.GetRoleByIds(req.Roles)
		if err != nil {
			return nil, fmt.Errorf("role with ids '%v' does not exist", req.Roles)
		}
		// Validar se o criador possui as roles necessárias ou é superusuário
		if !editor.IsSuperUser {
			if !utils.ContainsAll(editor.Roles, roles) {
				return nil, fmt.Errorf("failed to update user: editor does not have all required roles")
			}
		}

		// Atualizar as roles do usuário
		if err := s.GormStore.DB.Model(&user).Association("Roles").Replace(roles); err != nil {
			return nil, fmt.Errorf("failed to set roles for user: %v", err)
		}
	}

	// Atualizar outros campos do usuário
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Username = req.Username
	user.Email = req.Email
	user.Active = req.Active
	user.IsSuperUser = func() bool {
		if req.IsSuperUser {
			if editor.IsSuperUser {
				return true
			}
			panic(fmt.Errorf("only superusers can update other superusers"))
		}
		return false
	}()
	user.Phone1 = req.Phone1
	user.Phone2 = req.Phone2

	// Salvar as alterações
	if err := s.GormStore.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %s", err.Error())
	}
	return user, nil
}

func (s *Service) CreateRole(req schemas.CreateRoleRequest) (*models.Role, error) {
	// Verificar se o name já existe
	if err := s.GormStore.CheckIfRoleExistsByName(req.Name); err == nil {
		return nil, fmt.Errorf("role with name '%s' already exists", req.Name)
	}
	permissions, err := s.GormStore.GetPermissionByIds(req.Permissions)
	if err != nil {
		return nil, fmt.Errorf("permission with ids '%v' does not exist", req.Permissions)
	}

	var role models.Role

	role.Name = req.Name
	role.Permissions = permissions // Associar permissões à role

	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.GormStore.DB.Create(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %s", err.Error())
	}
	return &role, nil
}

func (s *Service) ListUsers() ([]models.User, error) {
	users, err := s.GormStore.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %s", err.Error())
	}
	return users, nil
}

func (s *Service) ListPermissions() ([]models.Permission, error) {
	permissions, err := s.GormStore.GetPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %s", err.Error())
	}
	return permissions, nil
}

func (s *Service) ListRoles() ([]models.Role, error) {
	roles, err := s.GormStore.GetRoles()
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %s", err.Error())
	}
	return roles, nil
}
