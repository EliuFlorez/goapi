package controller

import (
	"goapi/model"
	"goapi/request"
	"goapi/utils"
	"goapi/validator"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func EmailReset(context *gin.Context) {
	// Current User
	user, err := utils.CurrentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Token Error
	token, err := utils.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update Token
	err = user.UpdateUserByToken("email", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mail Send
	//utils.EmailSend(4440054, "Email Change", user.Email, user.FirstName, token)

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func EmailChange(context *gin.Context) {
	var input request.EmailResetInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.EmailResetForm{
		Email:             input.Email,
		EmailConfirmation: input.EmailConfirmation,
		Token:             input.Token,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Validate Token
	user, err := model.FindUserByToken("email_token", input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Equals
	if input.Email != input.EmailConfirmation {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email do not match"})
		return
	}

	// Change
	user.Email = input.Email

	// Clean Token
	user.EmailToken = ""
	user.EmailAt = time.Now()

	// Update
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"success": true})
}

func EmailToken(context *gin.Context) {
	token := context.Param("token")

	// Validate Token
	user, err := model.FindUserByColumn("email_token", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": user.ID > 0})
}

func EmailValid(context *gin.Context) {
	var input request.EmailInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.EmailForm{
		Email: input.Email,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Validate Token
	user, err := model.FindUserByEmail(input.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": user.ID > 0})
}
