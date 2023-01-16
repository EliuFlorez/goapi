package controller

import (
	"goapi/model"
	"goapi/request"
	"goapi/utils"
	"goapi/validator"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func AllPermissions(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "permission.all")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Find
	results, err := model.ListPermissions(10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pagination
	var pagination utils.Pagination
	pagination.Limit = 10
	pagination.TotalResults = model.CountPermissions()
	totalPages := int(math.Ceil(float64(pagination.TotalResults) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	pagination.Results = results

	context.JSON(http.StatusOK, &pagination)
}

func GetPermission(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "permission.show")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID
	id, _ := strconv.Atoi(context.Param("id"))

	// First
	role, err := model.GetPermission(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &role)
}

func AddPermission(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "permission.create")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.PermissionInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.PermissionForm{
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
	err = model.CreatePermission(input.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"success": true})
}

func UpdatePermission(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "permission.edit")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.PermissionInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.PermissionForm{
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
	id, _ := strconv.Atoi(context.Param("id"))

	// Edit
	err = model.EditPermission(id, input.Name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func DeletePermission(context *gin.Context) {
	// Validate Permission
	err := utils.ValidatePermission(context, "permission.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID
	id, _ := strconv.Atoi(context.Param("id"))

	// First
	err = model.DestroyPermission(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}
