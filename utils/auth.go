package utils

import (
	"errors"
	"goapi/model"

	"github.com/gin-gonic/gin"
)

func ValidateUserByPermission(context *gin.Context, permName string) (model.User, error) {
	// Current User
	user, err := CurrentUser(context)
	if err != nil {
		return model.User{}, err
	}

	// Check Permission
	ok, errRole := model.CheckPermission(user.ID, permName)
	if errRole != nil {
		return model.User{}, err
	}

	// Not Permission
	if !ok {
		return model.User{}, errors.New("user not permission")
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
	ok, errRole := model.CheckPermission(user.ID, permName)
	if errRole != nil {
		return err
	}

	// Not Permission
	if !ok {
		return errors.New("not permission")
	}

	return nil
}
