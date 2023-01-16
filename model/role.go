package model

import (
	"errors"
	"goapi/database"

	"gorm.io/gorm"
)

// Role represents the database model of roles
type Role struct {
	gorm.Model

	// Columns
	Name string
}

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TableName sets the table name
func (r Role) TableName() string {
	return "authority_roles"
}

// Response JSON
func (role *Role) ValueJson() (roleResponse *RoleResponse) {
	valueJson := &RoleResponse{
		ID:   role.ID,
		Name: role.Name,
	}

	return valueJson
}

// Count Roles
func CountRoles() int64 {
	var count int64

	database.Sql.Debug().Model(&Role{}).Count(&count)

	return count
}

// GetRoles returns all stored roles
func ListRoles(limit int) ([]Role, error) {
	var roles []Role
	err := database.Sql.Debug().Limit(limit).Model(&Role{}).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRole return first stored roles
func GetRole(id int) (Role, error) {
	var role Role
	err := database.Sql.Debug().Where("ID = ?", id).First(&role).Error
	if err != nil {
		return Role{}, err
	}

	return role, nil
}

// CreateRole stores a role in the database
// it accepts the role name. it returns an error
// in case of any
func CreateRole(roleName string) error {
	var dbRole Role
	res := database.Sql.Debug().Where("name = ?", roleName).First(&dbRole)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// create
			database.Sql.Create(&Role{Name: roleName})
			return nil
		}
	}

	return res.Error
}

// Edit Role
func EditRole(id int, name string) error {
	// check if the role
	var role Role
	res := database.Sql.Debug().Where("ID = ?", id).First(&role)
	if res.RowsAffected == 0 {
		return ErrRoleNotFound
	}

	// Update
	role.Name = name
	res = database.Sql.Debug().Model(&role).Updates(role)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// DeleteRole deletes a given role
// if the role is assigned to a user it returns an error
func DestroyRole(id int) error {
	// check if the role is assigned to a user
	var userRole UserRole
	res := database.Sql.Debug().Where("role_id = ?", id).First(&userRole)
	if res.Error == nil {
		return ErrRoleInUse
	}

	// revoke the assignment of permissions before deleting the role
	err := database.Sql.Debug().Where("role_id = ?", id).Delete(&RolePermission{}).Error
	if err != nil {
		return err
	}

	// delete the role
	err = database.Sql.Debug().Where("ID = ?", id).Delete(&Role{}).Error
	if err != nil {
		return err
	}

	return nil
}
