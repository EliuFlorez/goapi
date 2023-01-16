package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func StatusInput(context *gin.Context, input interface{}) {
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func StatusError(context *gin.Context, err error) {
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func ValidateFrom(context *gin.Context, formData interface{}) {
	// Validate Form
	v := validate.Struct(formData)
	v.StopOnError = false

	// Validate ?
	if !v.Validate() {
		context.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
		return
	}
}
