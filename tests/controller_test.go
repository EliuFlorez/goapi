package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goapi/controller"
	"goapi/database"
	"goapi/middleware"
	"goapi/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var body interface{}

// Un mapa para almacenar las sesiones activas.
var activeSessions = map[string]*Session{}

// Una estructura para representar una sesión.
type Session struct {
	Token     string
	CreatedAt time.Time
}

// Una función para crear una nueva sesión.
func CreateSession(token string) {
	session := &Session{
		Token:     token,
		CreatedAt: time.Now(),
	}
	activeSessions["0"] = session
}

// Una función para recuperar una sesión a partir de su ID.
func GetSession() *Session {
	session := activeSessions["0"]

	return session
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func router() *gin.Engine {
	router := gin.Default()

	// Auth
	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/signup", controller.Signup)
	//publicRoutes.POST("/signup/invite", controller.SignupInvite)
	//publicRoutes.GET("/signup/invite/:token", controller.SignupValid)
	publicRoutes.POST("/signin", controller.Signin)
	//publicRoutes.POST("/signin/code", controller.SigninCode)

	// Password
	//publicRoutes.GET("/password/token", controller.PasswordToken)
	publicRoutes.POST("/password/forgot", controller.PasswordForgot)
	//publicRoutes.POST("/password/reset", controller.PasswordReset)

	// Confirmation
	//publicRoutes.GET("/confirmation/token", controller.ConfirmationToken)
	//publicRoutes.POST("/confirmation/reset", controller.ConfirmationReset)

	// Email
	//publicRoutes.GET("/email/token", controller.EmailToken)
	//publicRoutes.POST("/email/change", controller.EmailChange)
	//publicRoutes.POST("/email/valid", controller.EmailValid)

	// Api - Protected
	protectedRoutes := router.Group("/api/v1")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())

	// Profile
	protectedRoutes.POST("/profile", controller.UpdateProfile)
	protectedRoutes.GET("/profile", controller.GetProfile)

	// Accounts
	protectedRoutes.POST("/accounts", controller.AddAccount)
	protectedRoutes.GET("/accounts", controller.GetAccounts)
	protectedRoutes.GET("/accounts/:id", controller.GetAccount)
	protectedRoutes.PUT("/accounts/:id", controller.UpdateAccount)
	//protectedRoutes.DELETE("/accounts/:id", controller.DeleteAccount)

	// Setting
	protectedRoutes.POST("/setting/email", controller.EmailReset)
	protectedRoutes.POST("/setting/password", controller.UserPassword)
	protectedRoutes.POST("/setting/twofa", controller.UserTwoFa)

	return router
}

func setup() {
	err := godotenv.Load("../.env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	database.Sql.AutoMigrate(&model.User{})
	database.Sql.AutoMigrate(&model.Account{})
	//database.Sql.AutoMigrate(&model.UserAccount{})
}

func teardown() {
	// migrator := database.Sql.Migrator()
	// migrator.DropTable(&model.UserAccount{})
	// migrator.DropTable(&model.User{})
	// migrator.DropTable(&model.Account{})
}

func makeRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+bearerToken())
	}
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	return response
}

func bearerToken() string {
	session := GetSession()

	return session.Token
}

func printResponse(response map[string]string) {
	fmt.Print("response-map: ")
	fmt.Println(response)
	fmt.Println("----")
}

func printHResponse(response *httptest.ResponseRecorder) {
	fmt.Print("response-http: ")
	fmt.Println(response)
	fmt.Println("----")
}
