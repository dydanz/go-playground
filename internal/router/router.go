package router

import (
	"go-playground/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	programRulesHandler *handler.ProgramRulesHandler,
	// ... other handlers ...
) *gin.Engine {
	router := gin.Default()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Program Rules routes
		programRulesGroup := v1.Group("/merchants")
		{
			programRulesGroup.GET("/:merchant_id/program-rules", programRulesHandler.GetProgramRulesByMerchantId)
			// ... other program rules routes ...
		}

		// ... other route groups ...
	}

	return router
}
