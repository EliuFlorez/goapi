package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetConnect(context *gin.Context) {
	// // Current User
	// user, err := utils.CurrentUser(context)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Type
	// typeParam := context.Param("type")
	// if typeParam == "" {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Oauth Type invalid."})
	// 	return
	// }

	// // Initialize hubspot client with OAuth refresh token.
	// client, _ := hubspot.NewClient(hubspot.SetOAuth(&hubspot.OAuthConfig{
	// 	GrantType:    hubspot.GrantTypeRefreshToken,
	// 	ClientID:     os.Getenv("HS_CLIENT_ID"),
	// 	ClientSecret: os.Getenv("HS_CLIENT_SECRET"),
	// 	RefreshToken: "YOUR_REFRESH_TOKEN",
	// }))

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func GetCallback(context *gin.Context) {
	// // Current User
	// user, err := utils.CurrentUser(context)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Type
	// typeParam := context.Param("type")
	// if typeParam == "" {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Oauth Type invalid."})
	// 	return
	// }

	// // Model
	// crm := model.Crm{
	// 	Name: typeParam,
	// }

	// // Save
	// crmNew, err := crm.Save()
	// if err != nil {
	// 	context.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	// 	return
	// }

	context.JSON(http.StatusOK, gin.H{"success": true})
}

func GetSelect(context *gin.Context) {
	// // Current User
	// user, err := utils.CurrentUser(context)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Type
	// typeParam := context.Param("type")
	// if typeParam == "" {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Oauth Type invalid."})
	// 	return
	// }

	// // ID
	// id, _ := strconv.Atoi(context.Param("id"))

	// // First
	// crm, err := model.FindCrmById(id)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	context.JSON(http.StatusOK, gin.H{"success": true})
}
