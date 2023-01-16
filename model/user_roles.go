package model

import (
	"errors"
	"goapi/database"

	"gorm.io/gorm"
)

// UserRole represents the relationship between users and roles
type UserRole struct {
	gorm.Model

	// Columns
	ID     uint
	UserID uint `gorm:"primarykey"`
	RoleID uint `gorm:"primarykey"`
}

// TableName sets the table name
func (u UserRole) TableName() string {
	return "authority_user_roles"
}

// GetUserRoles returns all user assigned roles
func GetUserRoles(userID uint) ([]Role, error) {
	var roles []Role
	var userRoles []UserRole
	database.Sql.Debug().Where("user_id = ?", userID).Find(&userRoles)

	for _, r := range userRoles {
		var role Role
		// for every user role get the role name
		res := database.Sql.Debug().Where("id = ?", r.RoleID).Find(&role)
		if res.Error == nil {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

// AssignRole assigns a given role to a user
// the first parameter is the user id, the second parameter is the role name
// if the role name doesn't have a matching record in the data base an error is returned
// if the user have already a role assigned to him an error is returned
func AssignRole(userID uint, roleName string) error {
	// make sure the role exist
	var role Role
	res := database.Sql.Debug().Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
	}

	// check if the role is already assigned
	var userRole UserRole
	res = database.Sql.Debug().Where("user_id = ?", userID).Where("role_id = ?", role.ID).First(&userRole)
	if res.Error == nil {
		//found a record, this role is already assigned to the same user
		return ErrRoleAlreadyAssigned
	}

	// assign the role
	err := database.Sql.Debug().Create(&UserRole{UserID: userID, RoleID: role.ID}).Error
	if err != nil {
		return err
	}

	return nil
}

// CheckRole checks if a role is assigned to a user
// it accepts the user id as the first parameter
// the role as the second parameter
// it returns an error if the role is not present in database
func CheckRole(userID uint, roleName string) (bool, error) {
	// find the role
	var role Role
	res := database.Sql.Debug().Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, ErrRoleNotFound
		}
	}

	// check if the role is a assigned
	var userRole UserRole
	res = database.Sql.Debug().Where("user_id = ?", userID).Where("role_id = ?", role.ID).First(&userRole)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}

	return true, nil
}

// RevokeRole revokes a user's role
// it returns a error in case of any
func RevokeRole(userID uint, roleName string) error {
	// find the role
	var role Role
	res := database.Sql.Debug().Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
	}

	// revoke the role
	err := database.Sql.Debug().Where("user_id = ?", userID).Where("role_id = ?", role.ID).Delete(UserRole{}).Error
	if err != nil {
		return err
	}

	return nil
}
