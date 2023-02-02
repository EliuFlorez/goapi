package main

import (
	"fmt"
	"goapi/controller"
	"goapi/database"
	"goapi/middleware"
	"goapi/model"
	"goapi/utils"
	"html"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	seedDatabase()
	serveApplication()
}

func loadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadDatabase() {
	// Connect
	database.Connect()

	// Create Table
	database.Sql.AutoMigrate(&model.User{})
	database.Sql.AutoMigrate(&model.Account{})
	database.Sql.AutoMigrate(&model.Crm{})

	// Roles & Permissions
	database.Sql.AutoMigrate(&model.Role{})
	database.Sql.AutoMigrate(&model.Permission{})
	database.Sql.AutoMigrate(&model.RolePermission{})
	database.Sql.AutoMigrate(&model.UserRole{})
}

func seedDatabase() {
	seedRoles := []string{"admin", "office", "owner", "user"}
	seedModules := []string{"roles", "permissions", "users", "accounts", "invoices", "crms", "prosperts"}
	seedPermissions := []string{"all", "show", "create", "edit", "destroy"}

	for _, roleName := range seedRoles {
		role := model.Role{Name: roleName, UserID: 0}
		database.Sql.Debug().FirstOrCreate(&role, role)
		if role.ID > 0 {
			if roleName == "user" {
				seedModules = []string{"crms", "prosperts"}
				seedPermissions = []string{"all", "show"}
			}
			for _, moduleName := range seedModules {
				for _, permissionName := range seedPermissions {
					modulePermission := moduleName + "." + permissionName
					permission := model.Permission{Name: modulePermission}
					database.Sql.Debug().FirstOrCreate(&permission, permission)
					if permission.ID > 0 {
						model.AssignPermissions(roleName, modulePermission)
					}
				}
			}
		}
	}

	hashedPassword, err := utils.HashPassword("demodemo")

	if err != nil {
		panic(err)
	}

	user := model.User{
		FirstName: "Admin",
		LastName:  "GO",
		Email:     html.EscapeString(strings.TrimSpace("admin@go.com")),
		Password:  hashedPassword,
	}

	database.Sql.Debug().FirstOrCreate(&user, model.User{Email: "admin@go.com"})

	if user.ID > 0 {
		model.AssignRole(user.ID, "admin")
	}
}

func serveApplication() {
	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Api status
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// CORS
	router.Use(middleware.Cors())

	// Auth
	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/signup", controller.Signup)
	publicRoutes.POST("/signup/invite", controller.SignupInvite)
	publicRoutes.GET("/signup/invite/:token", controller.SignupValid)
	publicRoutes.POST("/signin", controller.Signin)
	publicRoutes.POST("/signin/code", controller.SigninCode)
	// publicRoutes.GET("/logout", controller.Logout)

	// Password
	publicRoutes.GET("/password/token", controller.PasswordToken)
	publicRoutes.POST("/password/forgot", controller.PasswordForgot)
	publicRoutes.POST("/password/reset", controller.PasswordReset)

	// Confirmation
	publicRoutes.GET("/confirmation/token", controller.ConfirmationToken)
	publicRoutes.POST("/confirmation/reset", controller.ConfirmationReset)

	// Email
	publicRoutes.GET("/email/token", controller.EmailToken)
	publicRoutes.POST("/email/change", controller.EmailChange)
	publicRoutes.POST("/email/valid", controller.EmailValid)

	// Api - Protected
	protectedRoutes := router.Group("/api/v1")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())

	// Profile
	protectedRoutes.POST("/profile", controller.UpdateProfile)
	protectedRoutes.GET("/profile", controller.GetProfile)

	// Roles
	protectedRoutes.POST("/roles", controller.AddRole)
	protectedRoutes.GET("/roles", controller.AllRoles)
	protectedRoutes.GET("/roles/:id", controller.GetRole)
	protectedRoutes.PUT("/roles/:id", controller.UpdateRole)
	protectedRoutes.DELETE("/roles/:id", controller.DeleteRole)

	// Permissions
	protectedRoutes.POST("/permissions", controller.AddPermission)
	protectedRoutes.GET("/permissions", controller.AllPermissions)
	protectedRoutes.GET("/permissions/:id", controller.GetPermission)
	protectedRoutes.PUT("/permissions/:id", controller.UpdatePermission)
	protectedRoutes.DELETE("/permissions/:id", controller.DeletePermission)

	// Accounts
	protectedRoutes.POST("/accounts", controller.AddAccount)
	protectedRoutes.GET("/accounts", controller.GetAccounts)
	protectedRoutes.GET("/accounts/:id", controller.GetAccount)
	protectedRoutes.PUT("/accounts/:id", controller.UpdateAccount)
	protectedRoutes.DELETE("/accounts/:id", controller.DeleteAccount)

	// Setting
	protectedRoutes.POST("/setting/email", controller.EmailReset)
	protectedRoutes.POST("/setting/password", controller.UserPassword)
	protectedRoutes.POST("/setting/twofa", controller.UserTwoFa)

	// PORT
	router.Run(":8000")

	// TEST
	fmt.Println("Server running on port 8000")
}
