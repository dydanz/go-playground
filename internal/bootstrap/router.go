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
	UserHandler             *handler.UserHandler
	AuthHandler             *handler.AuthHandler
	PointsHandler           *handler.PointsHandler
	TransactionHandler      *handler.TransactionHandler
	RewardsHandler          *handler.RewardsHandler
	RedemptionHandler       *handler.RedemptionHandler
	PingHandler             *handler.PingHandler
	InternalLoadTestHandler *handler.InternalLoadTestHandler
	MerchantHandler         *handler.MerchantHandler
	ProgramHandler          *handler.ProgramHandler
	ProgramRuleHandler      *handler.ProgramRuleHandler
}

// InitializeHandlers initializes all handlers
func InitializeHandlers(services *Services, db *sql.DB, dbReplication *sql.DB, rdb *redislib.Client) *Handlers {
	return &Handlers{
		UserHandler:             handler.NewUserHandler(services.UserService),
		AuthHandler:             handler.NewAuthHandler(services.AuthService),
		PointsHandler:           handler.NewPointsHandler(services.PointsService),
		TransactionHandler:      handler.NewTransactionHandler(services.TransactionService),
		RewardsHandler:          handler.NewRewardsHandler(services.RewardsService),
		RedemptionHandler:       handler.NewRedemptionHandler(services.RedemptionService),
		PingHandler:             handler.NewPingHandler(db, dbReplication, rdb),
		InternalLoadTestHandler: handler.NewInternalLoadTestHandler(services.AuthService),
		MerchantHandler:         handler.NewMerchantHandler(services.MerchantService),
		ProgramHandler:          handler.NewProgramHandler(services.ProgramService),
		ProgramRuleHandler:      handler.NewProgramRuleHandler(services.ProgramRuleService),
	}
}

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(h *Handlers, authRepo *postgres.AuthRepository, sessionRepo redis.SessionRepository) *gin.Engine {
	r := gin.Default()

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
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
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
		}

		// Rewards routes
		rewards := api.Group("/rewards")
		{
			rewards.POST("", h.RewardsHandler.Create)
			rewards.GET("", h.RewardsHandler.GetAll)
			rewards.GET("/:id", h.RewardsHandler.GetByID)
			rewards.PUT("/:id", h.RewardsHandler.Update)
			rewards.DELETE("/:id", h.RewardsHandler.Delete)
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
			programRules.POST("", h.ProgramRuleHandler.Create)
			programRules.GET("/:id", h.ProgramRuleHandler.GetByID)
			programRules.GET("/program/:program_id", h.ProgramRuleHandler.GetByProgramID)
			programRules.PUT("/:id", h.ProgramRuleHandler.Update)
		}
	}

	// Protected HTML routes
	r.GET("/dashboard", middleware.AuthMiddleware(authRepo, sessionRepo), func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})

	// Add CSRF middleware
	r.Use(middleware.CSRFMiddleware())

	return r
}
