package controller

import (
	"goapi/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserPassword(context *gin.Context) {
	// User Current
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
	err = user.UpdateUserByToken("password", token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func UserTwoFa(context *gin.Context) {
	// User Current
	user, err := utils.CurrentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Status
	status := context.Param("status")

	// Update
	if status == "1" {
		user.SignInTwofa = true
	} else {
		user.SignInTwofa = false
	}

	// Update
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true})
}
