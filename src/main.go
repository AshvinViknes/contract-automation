package main

import (
	"smart-contract-automation/src/internal/controller"
	"smart-contract-automation/src/internal/repository"
	"smart-contract-automation/src/internal/router"
	"smart-contract-automation/src/internal/service"
	"smart-contract-automation/src/pkg/config"
	"smart-contract-automation/src/pkg/database"
	"smart-contract-automation/src/pkg/logger"
	"smart-contract-automation/src/pkg/server"

	"go.uber.org/zap"
)

func main() {
	// Load configuration with defaults
	cfg := config.New()

	// Initialize logger
	log := logger.Initialize(cfg.LogLevel, cfg.IsDevelopment())
	defer logger.Sync()

	// Set global logger
	zap.ReplaceGlobals(log)

	// Initialize database connection
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize layers
	contractRepo := repository.NewContractRepository(db)
	userRepo := repository.NewUserRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	billingRepo := repository.NewBillingRepository(db)

	contractService := service.NewContractService(contractRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, authService)
	workflowService := service.NewWorkflowService(workflowRepo)
	billingService := service.NewBillingService(billingRepo)

	contractController := controller.NewContractController(contractService)
	authController := controller.NewAuthController(authService)
	adminController := controller.NewAdminController(userService)
	setupController := controller.NewSetupController(userService, cfg)
	workflowController := controller.NewWorkflowController(workflowService)
	billingController := controller.NewBillingController(billingService)

	// Setup router
	r := router.SetupRouter(contractController, authController, adminController,
		setupController, workflowController, billingController)

	// Create and start server
	srv := server.NewServer(r, cfg.ServerPort)
	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
