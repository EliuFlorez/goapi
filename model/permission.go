package model

import (
	"errors"
	"goapi/database"

	"gorm.io/gorm"
)

// Permission represents the database model of permissions
type Permission struct {
	gorm.Model

	// Columns
	Name string `gorm:"index:idx_name"`
}

// TableName sets the table name
func (p Permission) TableName() string {
	return "authority_permissions"
}

// Count Permissions
func CountPermissions() int64 {
	var count int64

	database.Sql.Debug().Model(&Permission{}).Count(&count)

	return count
}

// GetPermissions returns all stored permissions
func ListPermissions(limit int) ([]Permission, error) {
	var permissions []Permission
	err := database.Sql.Debug().Limit(limit).Model(&Permission{}).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermission return first stored permission
func GetPermission(id int) (Permission, error) {
	var permission Permission
	err := database.Sql.Debug().Where("ID = ?", id).First(&permission).Error
	if err != nil {
		return Permission{}, err
	}

	return permission, nil
}

// CreatePermission stores a permission in the database
// it accepts the permission name. it returns an error
// in case of any
func CreatePermission(permName string) error {
	res := database.Sql.Debug().Where("name = ?", permName).First(&Permission{})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// create
			database.Sql.Debug().Create(&Permission{Name: permName})
			return nil
		}
	}

	return res.Error
}

// DeleteRole deletes a given role
// if the role is assigned to a user it returns an error
func EditPermission(id int, name string) error {
	// check if the permission
	var permission Permission
	res := database.Sql.Debug().Where("ID = ?", id).First(&permission)
	if res.RowsAffected == 0 {
		return ErrPermissionNotFound
	}

	// Update
	permission.Name = name
	res = database.Sql.Debug().Model(&permission).Updates(permission)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// DeletePermission deletes a given permission
// if the permission is assigned to a role it returns an error
func DestroyPermission(id int) error {
	// find the permission
	var permission Permission
	res := database.Sql.Debug().Where("ID = ?", id).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
	}

	// check if the permission is assigned to a role
	var rolePermission RolePermission
	res = database.Sql.Debug().Where("permission_id = ?", id).First(&rolePermission)
	if res.Error == nil {
		// role is assigned
		return ErrPermissionInUse
	}

	// delete the permission
	err := database.Sql.Debug().Where("ID = ?", id).Delete(&Permission{}).Error
	if err != nil {
		return err
	}

	return nil
}
