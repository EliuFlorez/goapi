package controller

import (
	"goapi/model"
	"goapi/request"
	"goapi/utils"
	"goapi/validator"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func AllRoles(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.all")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Find
	results, err := model.ListRoles(10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pagination
	var pagination utils.Pagination
	pagination.Limit = 10
	pagination.TotalResults = model.CountRoles()
	totalPages := int(math.Ceil(float64(pagination.TotalResults) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	pagination.Results = results

	context.JSON(http.StatusOK, &pagination)
}

func GetRole(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.show")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID
	id := utils.StringToUint(context.Param("id"))

	// First
	role, err := model.GetRole(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &role)
}

func AddRole(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.create")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.RoleInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.RoleForm{
		Name: input.Name,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Save
	err = model.CreateRole(input.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"success": true})
}

func UpdateRole(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.edit")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.RoleInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.RoleForm{
		Name: input.Name,
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

	// Edit
	err = model.EditRole(id, input.Name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteRole(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "roles.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID
	id := utils.StringToUint(context.Param("id"))

	// First
	err = model.DestroyRole(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}
