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

func PasswordForgot(context *gin.Context) {
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

	// Validate Email
	user, err := model.FindUserByEmail(input.Email)
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
	err = user.UpdateUserByToken("password", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mail Send
	//utils.EmailSend(4440054, "Password Forgot", user.Email, user.FirstName, token)

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func PasswordReset(context *gin.Context) {
	var input request.PasswordResetInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.PasswordResetForm{
		Password:             input.Password,
		PasswordConfirmation: input.PasswordConfirmation,
		Token:                input.Token,
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
	user, err := model.FindUserByToken("password_token", input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Password
	if input.Password != input.PasswordConfirmation {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Change
	user.Password = input.Password

	// Clean Token
	user.PasswordToken = ""
	user.PasswordAt = time.Now()

	// Update
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"success": true})
}

func PasswordToken(context *gin.Context) {
	// Token
	token := context.Param("token")

	// Validate Token
	user, err := model.FindUserByColumn("password_token", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": user.ID > 0})
}
