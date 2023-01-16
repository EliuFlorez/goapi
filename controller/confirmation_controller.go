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

func ConfirmationReset(context *gin.Context) {
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

	// Token Error
	token, err := utils.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update Token
	err = user.UpdateUserByToken("confirmation", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mail Send
	//utils.EmailSend(4440054, "Confirmation Email", user.Email, user.FirstName, token)

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func ConfirmationToken(context *gin.Context) {
	token := context.Param("token")

	// Validate Token
	user, err := model.FindUserByColumn("confirmation_token", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clean Token
	user.ConfirmationToken = ""
	user.ConfirmationAt = time.Now()
	user.ConfirmationEmail = true

	// Update
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": user.ID > 0})
}
