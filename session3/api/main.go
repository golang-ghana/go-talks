package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gorm-sqlc/config"
	"gorm-sqlc/database"
	gormAdvanced "gorm-sqlc/handlers/gorm"
	gormHandlers "gorm-sqlc/handlers/gorm"
	sqlcHandlers "gorm-sqlc/handlers/sqlc_"
	sqlcdb "gorm-sqlc/sqlc"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize GORM database
	gormDB, err := database.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to GORM database: %v", err)
	}

	// Auto-migrate disabled - using schema.sql for table creation

	// Initialize SQLC database
	sqlcDB, err := database.NewSQLCDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to SQLC database: %v", err)
	}
	defer sqlcDB.Close()

	// Initialize handlers
	gormUserHandler := gormHandlers.NewGormUserHandler(gormDB)
	gormTransactionHandler := gormHandlers.NewGormTransactionHandler(gormDB)
	gormAdvancedHandler := gormAdvanced.NewGormAdvancedHandler(gormDB)

	// Initialize SQLC queries with logging
	baseQueries := sqlcdb.New(sqlcDB)
	loggingQueries := database.NewLoggingQueries(baseQueries)
	sqlcUserHandler := sqlcHandlers.NewSQLCUserHandler(loggingQueries)
	sqlcTransactionHandler := sqlcHandlers.NewSQLCTransactionHandler(sqlcDB, baseQueries)

	// Setup router
	router := setupRouter(gormUserHandler, gormTransactionHandler, gormAdvancedHandler, sqlcUserHandler, sqlcTransactionHandler)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(
	gormUserHandler *gormHandlers.GormUserHandler,
	gormTransactionHandler *gormHandlers.GormTransactionHandler,
	gormAdvancedHandler *gormAdvanced.GormAdvancedHandler,
	sqlcUserHandler *sqlcHandlers.SQLCUserHandler,
	sqlcTransactionHandler *sqlcHandlers.SQLCTransactionHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// GORM API routes
	gormAPI := router.Group("/api/gorm")
	{
		// User CRUD operations
		gormUsers := gormAPI.Group("/users")
		{
			gormUsers.POST("", gormUserHandler.CreateUser)                   // CREATE
			gormUsers.GET("/:id", gormUserHandler.GetUserByID)               // READ by ID
			gormUsers.GET("/email/:email", gormUserHandler.GetUserByEmail)   // READ by email
			gormUsers.GET("", gormUserHandler.ListUsers)                     // LIST (paginated)
			gormUsers.GET("/active", gormUserHandler.ListActiveUsers)        // LIST active
			gormUsers.PUT("/:id", gormUserHandler.UpdateUser)                // UPDATE
			gormUsers.PATCH("/:id/status", gormUserHandler.UpdateUserStatus) // UPDATE status
			gormUsers.DELETE("/:id", gormUserHandler.DeleteUser)             // DELETE
			gormUsers.GET("/search", gormUserHandler.SearchUsersByName)      // SEARCH
			gormUsers.GET("/age-range", gormUserHandler.GetUsersByAgeRange)  // FILTER by age
			gormUsers.GET("/count", gormUserHandler.CountUsers)              // COUNT
		}

		// Advanced GORM operations
		gormUsers.GET("/complex", gormAdvancedHandler.ComplexQuery)         // Complex query
		gormUsers.GET("/partial", gormAdvancedHandler.SelectSpecificFields) // Select specific fields
		gormUsers.GET("/raw", gormAdvancedHandler.RawSQLQuery)              // Raw SQL

		// Transaction operations
		gormTransactions := gormAPI.Group("/transactions")
		{
			gormTransactions.POST("/users", gormTransactionHandler.CreateUsersInTransaction)
		}
	}

	// SQLC API routes
	sqlcAPI := router.Group("/api/sqlc")
	{
		// User CRUD operations
		sqlcUsers := sqlcAPI.Group("/users")
		{
			sqlcUsers.POST("", sqlcUserHandler.CreateUser)                   // CREATE
			sqlcUsers.GET("/:id", sqlcUserHandler.GetUserByID)               // READ by ID
			sqlcUsers.GET("/email/:email", sqlcUserHandler.GetUserByEmail)   // READ by email
			sqlcUsers.GET("", sqlcUserHandler.ListUsers)                     // LIST (paginated)
			sqlcUsers.GET("/active", sqlcUserHandler.ListActiveUsers)        // LIST active
			sqlcUsers.PUT("/:id", sqlcUserHandler.UpdateUser)                // UPDATE
			sqlcUsers.PATCH("/:id/status", sqlcUserHandler.UpdateUserStatus) // UPDATE status
			sqlcUsers.DELETE("/:id", sqlcUserHandler.DeleteUser)             // DELETE
			sqlcUsers.GET("/search", sqlcUserHandler.SearchUsersByName)      // SEARCH
			sqlcUsers.GET("/age-range", sqlcUserHandler.GetUsersByAgeRange)  // FILTER by age
			sqlcUsers.GET("/count", sqlcUserHandler.CountUsers)              // COUNT
		}

		// Transaction operations
		sqlcTransactions := sqlcAPI.Group("/transactions")
		{
			sqlcTransactions.POST("/users", sqlcTransactionHandler.CreateUsersInTransaction)
		}
	}

	return router
}
