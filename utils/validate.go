package utils

import (
	"errors"
	"goapi/model"

	"github.com/gin-gonic/gin"
)

func ValidateRole(context *gin.Context, roleName string) error {
	// Current User
	user, err := CurrentUser(context)
	if err != nil {
		return err
	}

	// Check Role
	ok, errRole := model.CheckRole(user.ID, roleName)
	if errRole != nil {
		return err
	}

	// Not Role
	if !ok {
		return errors.New("not permission")
	}

	return nil
}

func ValidateUserByRole(context *gin.Context, roleName string) (model.User, error) {
	// Current User
	user, err := CurrentUser(context)
	if err != nil {
		return model.User{}, err
	}

	// Check Role
	ok, errRole := model.CheckRole(user.ID, roleName)
	if errRole != nil {
		return model.User{}, err
	}

	// Not Role
	if !ok {
		return model.User{}, errors.New("user not role")
	}

	return user, nil
}

func ValidatePermission(context *gin.Context, permName string) error {
	// Current User
	user, err := CurrentUser(context)
	if err != nil {
		return err
	}

	// Check Permission
	ok, errPermission := model.CheckPermission(user.ID, permName)
	if errPermission != nil {
		return err
	}

	// Not Permission
	if !ok {
		return errors.New("not permission")
	}

	return nil
}

func ValidateUserByPermission(context *gin.Context, permName string) (model.User, error) {
	// Current User
	user, err := CurrentUser(context)
	if err != nil {
		return model.User{}, err
	}

	// Check Permission
	ok, errPermission := model.CheckPermission(user.ID, permName)
	if errPermission != nil {
		return model.User{}, err
	}

	// Not Permission
	if !ok {
		return model.User{}, errors.New("user not permission")
	}

	return user, nil
}
