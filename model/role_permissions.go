package model

import (
	"errors"
	"goapi/database"

	"gorm.io/gorm"
)

// RolePermission stores the relationship between roles and permissions
type RolePermission struct {
	gorm.Model

	// Columns
	ID           uint
	RoleID       uint `gorm:"primarykey"`
	PermissionID uint `gorm:"primarykey"`
}

// TableName sets the table name
func (r RolePermission) TableName() string {
	return "authority_role_permissions"
}

// AssignPermissions assigns a group of permissions to a given role
// it accepts in the first parameter the role name, it returns an error if there is not matching record
// of the role name in the database.
// the second parameter is a slice of strings which represents a group of permissions to be assigned to the role
// if any of these permissions doesn't have a matching record in the database the operations stops, changes reverted
// and error is returned
// in case of success nothing is returned
func AssignPermissions(roleName string, permName string) error {
	// get the role id
	var role Role
	roleData := database.Sql.Debug().Where("name = ?", roleName).First(&role)
	if roleData.Error != nil {
		if errors.Is(roleData.Error, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
	}

	var perm Permission

	// get the permissions ids
	permissionData := database.Sql.Debug().Where("name = ?", permName).First(&perm)
	if permissionData.Error != nil {
		if errors.Is(permissionData.Error, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
	}

	// ignore any assigned permission
	var rolePerm RolePermission
	res := database.Sql.Debug().Where("role_id = ?", role.ID).Where("permission_id =?", perm.ID).First(&rolePerm)
	if res.Error != nil {
		// assign the record
		cRes := database.Sql.Debug().Create(&RolePermission{RoleID: role.ID, PermissionID: perm.ID})
		if cRes.Error != nil {
			return cRes.Error
		}
	}

	return nil
}

// CheckPermission checks if a permission is assigned to the role that's assigned to the user.
// it accepts the user id as the first parameter
// the permission as the second parameter
// it returns an error if the permission is not present in the database
func CheckPermission(userID uint, permName string) (bool, error) {
	// the user role
	var userRoles []UserRole
	res := database.Sql.Debug().Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}

	//prepare an array of role ids
	var roleIDs []uint
	for _, r := range userRoles {
		roleIDs = append(roleIDs, r.RoleID)
	}

	// find the permission
	var perm Permission
	res = database.Sql.Debug().Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, ErrPermissionNotFound
		}
	}

	// find the role permission
	var rolePermission RolePermission
	res = database.Sql.Debug().Where("role_id IN (?)", roleIDs).Where("permission_id = ?", perm.ID).First(&rolePermission)
	if res.Error != nil {
		return false, nil
	}

	return true, nil
}

// CheckRolePermission checks if a role has the permission assigned
// it accepts the role as the first parameter
// it accepts the permission as the second parameter
// it returns an error if the role is not present in database
// it returns an error if the permission is not present in database
func CheckRolePermission(roleName string, permName string) (bool, error) {
	// find the role
	var role Role
	res := database.Sql.Debug().Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, ErrRoleNotFound
		}
	}

	// find the permission
	var perm Permission
	res = database.Sql.Debug().Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, ErrPermissionNotFound
		}
	}

	// find the rolePermission
	var rolePermission RolePermission
	res = database.Sql.Debug().Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).First(&rolePermission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}

	return true, nil
}

// RevokePermission revokes a permission from the user's assigned role
// it returns an error in case of any
func RevokePermission(userID uint, permName string) error {
	// revoke the permission from all roles of the user
	// find the user roles
	var userRoles []UserRole
	res := database.Sql.Debug().Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
	}

	// find the permission
	var perm Permission
	res = database.Sql.Debug().Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
	}

	// revoke the permission
	for _, r := range userRoles {
		database.Sql.Debug().Where("role_id = ?", r.RoleID).Where("permission_id = ?", perm.ID).Delete(RolePermission{})
	}

	return nil
}

// RevokeRolePermission revokes a permission from a given role
// it returns an error in case of any
func RevokeRolePermission(roleName string, permName string) error {
	// find the role
	var role Role
	res := database.Sql.Debug().Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
	}

	// find the permission
	var perm Permission
	res = database.Sql.Debug().Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
	}

	// revoke the permission
	err := database.Sql.Debug().Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).Delete(RolePermission{}).Error
	if err != nil {
		return err
	}

	return nil
}
