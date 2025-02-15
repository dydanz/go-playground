package bootstrap

import (
	"go-playground/internal/handler"
	"go-playground/internal/middleware"
	"go-playground/internal/repository/postgres"
	"go-playground/internal/repository/redis"
	"go-playground/internal/util"
	"net/http"

	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	redislib "github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers holds all handler instances
type Handlers struct {
	UserHandler              *handler.UserHandler
	AuthHandler              *handler.AuthHandler
	PointsHandler            *handler.PointsHandler
	TransactionHandler       *handler.TransactionHandler
	RewardsHandler           *handler.RewardsHandler
	RedemptionHandler        *handler.RedemptionHandler
	PingHandler              *handler.PingHandler
	InternalLoadTestHandler  *handler.InternalLoadTestHandler
	MerchantHandler          *handler.MerchantHandler
	MerchantCustomersHandler *handler.MerchantCustomersHandler
	ProgramHandler           *handler.ProgramHandler
	ProgramRulesHandler      *handler.ProgramRulesHandler
}

// InitializeHandlers initializes all handlers
func InitializeHandlers(services *Services, db *sql.DB, dbReplication *sql.DB, rdb *redislib.Client) *Handlers {
	return &Handlers{
		UserHandler:              handler.NewUserHandler(services.UserService),
		AuthHandler:              handler.NewAuthHandler(services.AuthService),
		PointsHandler:            handler.NewPointsHandler(services.PointsService),
		TransactionHandler:       handler.NewTransactionHandler(services.TransactionService),
		RewardsHandler:           handler.NewRewardsHandler(services.RewardsService),
		RedemptionHandler:        handler.NewRedemptionHandler(services.RedemptionService),
		PingHandler:              handler.NewPingHandler(db, dbReplication, rdb),
		InternalLoadTestHandler:  handler.NewInternalLoadTestHandler(services.AuthService),
		MerchantHandler:          handler.NewMerchantHandler(services.MerchantService),
		MerchantCustomersHandler: handler.NewMerchantCustomersHandler(services.MerchantCustomersService),
		ProgramHandler:           handler.NewProgramHandler(services.ProgramService),
		ProgramRulesHandler:      handler.NewProgramRulesHandler(services.ProgramRuleService),
	}
}

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(h *Handlers, authRepo *postgres.AuthRepository, sessionRepo redis.SessionRepository) *gin.Engine {
	r := gin.Default()

	// Debug mode
	gin.SetMode(gin.DebugMode)

	// Load templates explicitly
	r.LoadHTMLFiles(
		"internal/static/pages/sign-in.html",
		"internal/static/pages/sign-up.html",
		"internal/static/pages/dashboard.html",
		"internal/static/pages/profile.html",
		"internal/static/pages/billing.html",
		"internal/static/pages/transactions.html",
		"internal/static/pages/merchants.html",
		"internal/static/pages/components/navbar.tmpl",
		"internal/static/pages/components/sidenav.tmpl",
		"internal/static/pages/components/sidenav-card.tmpl",
		"internal/static/pages/components/add-merchant-modal.html",
	)

	// Serve static files
	r.Static("/static/pages", "internal/static/pages")
	r.Static("/static/assets", "internal/static/assets")

	// Add latency middleware globally
	r.Use(util.GinHandlerLatencyDecorator())

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-Id"}
	r.Use(cors.New(config))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	r.GET("/ping", h.PingHandler.Ping)

	r.GET("/sign-in", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sign-in.html", nil)
	})
	r.GET("/sign-up", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sign-up.html", nil)
	})

	// Public auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.AuthHandler.Register)
		auth.POST("/verify", h.AuthHandler.Verify)
		auth.POST("/login", h.AuthHandler.Login)

		// FOR LOAD TEST ONLY
		auth.GET("/test/get-verification/code", h.InternalLoadTestHandler.GetVerificationCode)
		auth.GET("/test/random-user", h.InternalLoadTestHandler.GetRandomVerifiedUser)
	}

	// Protected routes with auth middleware
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authRepo, sessionRepo))
	{
		api.POST("/auth/logout", h.AuthHandler.Logout)

		// Users routes
		users := api.Group("/users")
		{
			users.GET("/me", h.UserHandler.GetMe)
			users.GET("", h.UserHandler.GetAll)
			users.GET("/:id", h.UserHandler.GetByID)
			users.POST("", h.UserHandler.Create)
			users.PUT("/:id", h.UserHandler.Update)
			users.DELETE("/:id", h.UserHandler.Delete)
		}

		// Points routes
		points := api.Group("/points")
		{
			points.GET("/:customer_id/:program_id/ledger", h.PointsHandler.GetLedger)
			points.GET("/:customer_id/:program_id/balance", h.PointsHandler.GetBalance)
			points.POST("/:customer_id/:program_id/earn", h.PointsHandler.EarnPoints)
			points.POST("/:customer_id/:program_id/redeem", h.PointsHandler.RedeemPoints)
		}

		// Transactions routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("", h.TransactionHandler.Create)
			transactions.GET("/:id", h.TransactionHandler.GetByID)
			transactions.GET("/user/:user_id", h.TransactionHandler.GetByCustomerID)
			transactions.GET("/merchant/:merchant_id", h.TransactionHandler.GetByMerchantID)
		}

		// Rewards routes
		rewards := api.Group("/rewards")
		{
			rewards.POST("", h.RewardsHandler.Create)
			rewards.GET("", h.RewardsHandler.GetAll)
			rewards.GET("/:id", h.RewardsHandler.GetByID)
			rewards.PUT("/:id", h.RewardsHandler.Update)
			rewards.DELETE("/:id", h.RewardsHandler.Delete)
			rewards.GET("/program/:program_id", h.RewardsHandler.GetByProgramID)
		}

		// Redemptions routes
		redemptions := api.Group("/redemptions")
		{
			redemptions.POST("", h.RedemptionHandler.Create)
			redemptions.GET("/:id", h.RedemptionHandler.GetByID)
			redemptions.GET("/user/:user_id", h.RedemptionHandler.GetByUserID)
			redemptions.PUT("/:id/status", h.RedemptionHandler.UpdateStatus)
		}

		// Merchants routes
		merchants := api.Group("/merchants")
		{
			merchants.POST("", h.MerchantHandler.Create)
			merchants.GET("", h.MerchantHandler.GetAll)
			merchants.GET("/:id", h.MerchantHandler.GetByID)
			merchants.PUT("/:id", h.MerchantHandler.Update)
			merchants.DELETE("/:id", h.MerchantHandler.Delete)
			merchants.GET("/user/:user_id", h.MerchantHandler.GetMerchantsByUserID)
		}

		// Merchant Customers routes
		merchantCustomers := api.Group("/merchant-customers")
		{
			merchantCustomers.POST("", h.MerchantCustomersHandler.Create)
			merchantCustomers.GET("/:id", h.MerchantCustomersHandler.GetByID)
			merchantCustomers.GET("/merchant/:merchant_id", h.MerchantCustomersHandler.GetByMerchantID)
			merchantCustomers.PUT("/:id", h.MerchantCustomersHandler.Update)
			merchantCustomers.POST("/login", h.MerchantCustomersHandler.ValidateCredentials)
		}

		// Programs routes
		programs := api.Group("/programs")
		{
			programs.POST("", h.ProgramHandler.Create)
			programs.GET("/:id", h.ProgramHandler.GetByID)
			programs.GET("/merchant/:merchant_id", h.ProgramHandler.GetByMerchantID)
			programs.PUT("/:id", h.ProgramHandler.Update)
			programs.DELETE("/:id", h.ProgramHandler.Delete)
		}

		programRules := api.Group("/program-rules")
		{
			programRules.POST("", h.ProgramRulesHandler.Create)
			programRules.GET("/:id", h.ProgramRulesHandler.GetByID)
			programRules.GET("/program/:program_id", h.ProgramRulesHandler.GetByProgramID)
			programRules.PUT("/:id", h.ProgramRulesHandler.Update)
		}
	}

	// Protected HTML routes
	r.GET("/dashboard", middleware.AuthMiddleware(authRepo, sessionRepo), middleware.CSRFMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})
	r.GET("/profile", middleware.AuthMiddleware(authRepo, sessionRepo), middleware.CSRFMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})
	r.GET("/transactions", middleware.AuthMiddleware(authRepo, sessionRepo), middleware.CSRFMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "transactions.html", nil)
	})
	r.GET("/merchants", middleware.AuthMiddleware(authRepo, sessionRepo), middleware.CSRFMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "merchants.html", nil)
	})
	r.GET("/billing", middleware.AuthMiddleware(authRepo, sessionRepo), middleware.CSRFMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "billing.html", nil)
	})

	// Add CSRF middleware
	r.Use(middleware.CSRFMiddleware())

	return r
}
