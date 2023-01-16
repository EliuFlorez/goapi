package controller

import (
	"goapi/model"
	"goapi/request"
	"goapi/utils"
	"goapi/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func GetProfile(context *gin.Context) {
	// Current User
	user, err := utils.CurrentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts, err := model.AllAccounts(&user, 10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response
	userResponse := &model.UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Accounts:  accounts,
	}

	context.JSON(http.StatusOK, &userResponse)
}

func UpdateProfile(context *gin.Context) {
	var input request.ProfileInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Current User
	user, err := utils.CurrentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.ProfileForm{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Edit
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Phone = input.Phone

	// Update
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}
