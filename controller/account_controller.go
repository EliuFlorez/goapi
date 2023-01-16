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

func GetAccounts(context *gin.Context) {
	// Validate User by permission
	user, err := utils.ValidateUserByPermission(context, "account.all")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model
	var account model.Account

	// Validate Find
	accounts, err := account.All(&user, 10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pagination
	var pagination utils.Pagination
	totalResults := account.Total(&user)
	pagination.TotalResults = totalResults
	totalPages := int(math.Ceil(float64(totalResults) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	pagination.Results = accounts

	context.JSON(http.StatusOK, &pagination)
}

func GetAccount(context *gin.Context) {
	// Validate User by permission
	user, err := utils.ValidateUserByPermission(context, "account.show")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model
	var account model.Account

	// ID
	id, _ := strconv.Atoi(context.Param("id"))

	// First
	accountGet, err := account.Get(&user, id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &accountGet)
}

func AddAccount(context *gin.Context) {
	// Validate User by permission
	user, err := utils.ValidateUserByPermission(context, "account.create")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.AccountInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.AccountForm{
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
	err = user.AddAccountToUser(input.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"success": true})
}

func UpdateAccount(context *gin.Context) {
	// Validate User by permission
	user, err := utils.ValidateUserByPermission(context, "account.edit")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input
	var input request.AccountInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.AccountForm{
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

	// Models
	var accountModel model.Account

	// ID
	id, _ := strconv.Atoi(context.Param("id"))

	// First
	account, err := accountModel.Get(&user, id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update
	account.Name = input.Name
	err = account.UpdateAccountByUser(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &account)
}

func DeleteAccount(context *gin.Context) {
	// Validate User by permission
	user, err := utils.ValidateUserByPermission(context, "account.destroy")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model
	var accountModel model.Account

	// ID
	id, _ := strconv.Atoi(context.Param("id"))

	// First
	account, err := accountModel.Get(&user, id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Delete
	err = account.DestroyAccountByUser(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}
