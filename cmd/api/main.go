package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	// Swagger docs
	_ "finanzas-backend/cmd/api/docs"

	// Shared
	"finanzas-backend/internal/shared/infrastructure/config"
	"finanzas-backend/internal/shared/infrastructure/persistence"

	// IAM
	iamACLImpl "finanzas-backend/internal/iam/application/acl"
	iamCommandServices "finanzas-backend/internal/iam/application/commandservices"
	iamQueryServices "finanzas-backend/internal/iam/application/queryservices"
	iamRepos "finanzas-backend/internal/iam/infrastructure/persistence/repositories"
	iamSecurity "finanzas-backend/internal/iam/infrastructure/security"
	iamACL "finanzas-backend/internal/iam/interfaces/acl"
	iamControllers "finanzas-backend/internal/iam/interfaces/rest/controllers"

	// Mortgage
	mortgageACL "finanzas-backend/internal/mortgage/application/acl"
	mortgageCommandServices "finanzas-backend/internal/mortgage/application/commandservices"
	mortgageQueryServices "finanzas-backend/internal/mortgage/application/queryservices"
	mortgageRepos "finanzas-backend/internal/mortgage/infrastructure/persistence/repositories"
	mortgageControllers "finanzas-backend/internal/mortgage/interfaces/rest/controllers"
	mortgageMiddleware "finanzas-backend/internal/mortgage/interfaces/rest/middleware"
)

// @title Finanzas API - MiVivienda Mortgage Calculator
// @version 1.0
// @description API for IAM and Mortgage (French Method) bounded contexts
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup Gin
	router := gin.Default()

	// Setup CORS
	router.Use(corsMiddleware())

	// Setup dependencies and routes
	iamFacade := setupIAMContext(router, db, cfg)
	setupMortgageContext(router, db, iamFacade)

	// Swagger UI route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root redirect to Swagger
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	swaggerURL := fmt.Sprintf("http://%s/swagger/index.html", serverAddr)

	fmt.Println("=====================================")
	fmt.Println("Server starting...")
	fmt.Println("=====================================")
	fmt.Printf("Server: http://%s\n", serverAddr)
	fmt.Printf("Swagger UI: %s\n", swaggerURL)
	fmt.Println("=====================================")

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware configura CORS para producción y desarrollo
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupIAMContext(router *gin.Engine, db *gorm.DB, cfg *config.Config) iamACL.IAMContextFacade {
	// JWT Service
	jwtService := iamSecurity.NewJWTService(
		cfg.JWT.SecretKey,
		cfg.JWT.Issuer,
		cfg.JWT.ExpirationHrs,
	)

	// Repositories
	userRepo := iamRepos.NewUserRepository(db)

	// Services
	userCommandService := iamCommandServices.NewUserCommandService(userRepo)
	userQueryService := iamQueryServices.NewUserQueryService(userRepo)
	authService := iamCommandServices.NewAuthenticationService(userRepo, jwtService)

	// ACL Facade (expuesto a otros bounded contexts)
	iamFacade := iamACLImpl.NewIAMContextFacade(jwtService, userRepo)

	// External Services (ACL for own middleware)
	externalAuthService := mortgageACL.NewExternalAuthenticationService(iamFacade)

	// Middleware
	authMiddleware := mortgageMiddleware.JWTAuthMiddleware(externalAuthService)

	// Controllers
	userController := iamControllers.NewUserController(userCommandService, userQueryService, authService)

	// Routes
	iamGroup := router.Group("/api/v1/iam")
	{
		iamGroup.POST("/register", userController.Register)
		iamGroup.POST("/login", userController.Login)

		// Protected routes
		iamGroup.PUT("/profile", authMiddleware, userController.UpdateProfile)
	}

	return iamFacade
}

func setupMortgageContext(router *gin.Engine, db *gorm.DB, iamFacade iamACL.IAMContextFacade) {
	// External Services (ACL)
	externalAuthService := mortgageACL.NewExternalAuthenticationService(iamFacade)

	// Middleware
	authMiddleware := mortgageMiddleware.JWTAuthMiddleware(externalAuthService)

	// Repositories
	mortgageRepo := mortgageRepos.NewMortgageRepository(db)

	// Services
	mortgageCommandService := mortgageCommandServices.NewMortgageCommandService(mortgageRepo)
	mortgageQueryService := mortgageQueryServices.NewMortgageQueryService(mortgageRepo)

	// Controllers
	mortgageController := mortgageControllers.NewMortgageController(mortgageCommandService, mortgageQueryService)

	// Routes (todas protegidas con JWT)
	mortgageGroup := router.Group("/api/v1/mortgage")
	mortgageGroup.Use(authMiddleware) // Aplicar middleware a todas las rutas
	{
		mortgageGroup.POST("/calculate", mortgageController.CalculateMortgage)
		mortgageGroup.GET("/:id", mortgageController.GetMortgageByID)
		mortgageGroup.PUT("/:id", mortgageController.UpdateMortgage)
		mortgageGroup.DELETE("/:id", mortgageController.DeleteMortgage)
		mortgageGroup.GET("/history", mortgageController.GetMortgageHistory)
	}
}
