package router

import (
	"net/http"
	"smart-contract-automation/src/internal/controller"
	"smart-contract-automation/src/pkg/middleware"

	"time"

	"smart-contract-automation/src/internal/model"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRouter(contractController *controller.ContractController, authController *controller.AuthController,
	adminController *controller.AdminController, setupController *controller.SetupController,
	workflowController *controller.WorkflowController, billingController *controller.BillingController) *chi.Mux {
	r := chi.NewRouter()

	// Core middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Post("/api/v1/auth/login", authController.Login)
		r.Post("/api/v1/setup/super-admin", setupController.CreateFirstSuperAdmin)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authController.AuthService))

		// Admin management routes (super_admin only)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(model.RoleSuperAdmin))
			r.Route("/api/v1/admins", func(r chi.Router) {
				r.Post("/", adminController.CreateAdmin)
				r.Get("/", adminController.ListAdmins)
				r.Get("/{id}", adminController.GetAdmin)
				r.Delete("/{id}", adminController.DeleteAdmin)
			})
		})

		r.Route("/api/v1/contracts", func(r chi.Router) {
			// Routes accessible by all authenticated users
			r.Get("/", contractController.ListContracts)
			r.Get("/{id}", contractController.GetContract)

			// Routes accessible by admin and super_admin
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole(model.RoleAdmin, model.RoleSuperAdmin))
				r.Post("/", contractController.CreateContract)
				r.Put("/{id}", contractController.UpdateContract)
				r.Get("/{id}/status", contractController.GetContractStatus)
			})

			// Routes accessible only by super_admin
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole(model.RoleSuperAdmin))
				r.Patch("/{id}/deploy", contractController.DeployContract)
				r.Patch("/{id}/terminate", contractController.TerminateContract)
				r.Delete("/{id}", contractController.DeleteContract)
			})
		})

		r.Route("/api/v1/analytics", func(r chi.Router) {
			r.Use(middleware.RequireRole(model.RoleAdmin, model.RoleSuperAdmin))
			r.Get("/contracts/status/count", contractController.GetContractStatusCounts)
			r.Get("/contracts", contractController.GetContractsByStatus)
			r.Get("/contracts/trend", contractController.GetContractTrend)
		})

		// Workflow routes
		r.Route("/api/v1/workflows", func(r chi.Router) {
			r.Use(middleware.RequireRole(model.RoleAdmin, model.RoleSuperAdmin))
			r.Post("/", workflowController.CreateWorkflow)
			r.Get("/", workflowController.ListWorkflows)
			r.Get("/{id}", workflowController.GetWorkflow)
			r.Put("/{id}/approve", workflowController.ApproveWorkflow)
			r.Put("/{id}/reject", workflowController.RejectWorkflow)
		})

		// Billing routes
		r.Route("/api/v1/billing", func(r chi.Router) {
			r.Use(middleware.RequireRole(model.RoleAdmin, model.RoleSuperAdmin))
			r.Post("/", billingController.CreateBill)
			r.Get("/{id}", billingController.GetBill)
			r.Put("/{id}", billingController.UpdateBill)
			r.Delete("/{id}", billingController.DeleteBill)
			r.Get("/", billingController.ListBills)
			r.Get("/user/{id}", billingController.ListBillsByUserID)
		})
	})

	return r
}
