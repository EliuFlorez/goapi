package controller

import (
	"goapi/model"
	"goapi/request"
	"goapi/utils"
	"goapi/validator"
	"html"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func GetUsers(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Find
	users, err := model.AllUsers(10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pagination
	var pagination utils.Pagination
	totalResults := model.TotalUsers()
	pagination.TotalResults = totalResults
	totalPages := int(math.Ceil(float64(totalResults) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	pagination.Results = users

	context.JSON(http.StatusOK, &pagination)
}

func GetUser(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID
	id := utils.StringToUint(context.Param("id"))

	// First
	userData, err := model.FindUserById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &userData)
}

func AddUser(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.UserInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.UserForm{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Model
	userNew := model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     html.EscapeString(strings.TrimSpace(input.Email)),
	}

	// Save
	newUser, err := userNew.Save()
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		context.JSON(http.StatusConflict, gin.H{"error": "User with that email already exists"})
		return
	} else if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": newUser.ID > 0})
}

func UpdateUser(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.UserInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.UserForm{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// ID
	id := utils.StringToUint(context.Param("id"))

	// First
	userData, err := model.FindUserById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Edit
	userData.FirstName = input.FirstName
	userData.LastName = input.LastName
	userData.Email = input.Email

	// Update
	err = userData.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteUser(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID
	id := utils.StringToUint(context.Param("id"))

	// First
	userData, err := model.FindUserById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Delete
	err = model.DeleteUserById(userData.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}
