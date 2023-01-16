package controller

import (
	"goapi/model"
	"goapi/request"
	"goapi/utils"
	"goapi/validator"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func Signup(context *gin.Context) {
	var input request.SignUpInput

	// Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := validator.SignUpForm{
		CompanyName:          input.CompanyName,
		FirstName:            input.FirstName,
		LastName:             input.LastName,
		Email:                input.Email,
		Password:             input.Password,
		PasswordConfirmation: input.PasswordConfirmation,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Password Match to confirmation
	if input.Password != input.PasswordConfirmation {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Hash Password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Model
	user := model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     html.EscapeString(strings.TrimSpace(input.Email)),
		Password:  hashedPassword,
	}

	// Save
	newUser, err := user.Save()
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		context.JSON(http.StatusConflict, gin.H{"error": "User with that email already exists"})
		return
	} else if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Add Account
	err = user.AddAccountToUser(input.CompanyName)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mail Send
	//utils.EmailSend(4440054, "Register", savedUser.Email, savedUser.FirstName, "Thanks")

	context.JSON(http.StatusOK, gin.H{"success": newUser.ID > 0})
}

func SignupInvite(context *gin.Context) {
	var input request.SignUpInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.SignUpForm{
		Password:             input.Password,
		PasswordConfirmation: input.PasswordConfirmation,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Password Match to confirmation
	if input.Password != input.PasswordConfirmation {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Token
	token := context.Param("token")

	// Validate Token
	user, err := model.FindUserByToken("invitation_token", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clean Token
	user.InvitationToken = ""
	user.InvitationAt = time.Now()

	// Save
	savedUser, err := user.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mail Send
	//utils.EmailSend(4440054, "Register", savedUser.Email, savedUser.FirstName, "Thanks")

	context.JSON(http.StatusOK, gin.H{"success": savedUser.ID > 0})
}

func SignupValid(context *gin.Context) {
	// Token
	token := context.Param("token")

	// Validate Token
	user, err := model.FindUserByColumn("invitation_token", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": user.ID > 0})
}

func Signin(context *gin.Context) {
	var input request.SignInInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.SignInForm{
		Email:    input.Email,
		Password: input.Password,
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

	// Validate Password
	err = utils.ValidatePassword(user.Password, input.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Two Fa
	if user.SignInTwofa {
		// Mail Send
		//utils.EmailSend(4440054, "Code Security", user.Email, user.FirstName, "Code")
		context.JSON(http.StatusCreated, gin.H{"success": true})
	} else {
		// Token Error
		token, err := utils.GenerateJWT(user)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func SigninCode(context *gin.Context) {
	var input request.SignInCodeInput

	// Validate Input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator
	formData := &validator.SignInCodeForm{
		Code: input.Code,
	}

	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}

	// Validate Code
	user, err := model.FindUserByCode(input.Code)
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

	// Code Clean
	user.TwofaCode = ""

	// Update
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}
