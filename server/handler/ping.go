package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type PingHandler struct {
	dbPrimary   *sql.DB
	dbReplica   *sql.DB
	redisClient *redis.Client
}

func NewPingHandler(dbPrimary *sql.DB, dbReplica *sql.DB, redisClient *redis.Client) *PingHandler {
	return &PingHandler{
		dbPrimary:   dbPrimary,
		dbReplica:   dbReplica,
		redisClient: redisClient,
	}
}

// @Summary Health check
// @Description Check the health of database and Redis connections
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ping [get]
func (h *PingHandler) Ping(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"services": map[string]interface{}{
			"database": map[string]interface{}{
				"primary": map[string]string{
					"status": "unknown",
				},
				"replica": map[string]string{
					"status": "unknown",
				},
			},
			"redis": map[string]string{
				"status": "unknown",
			},
		},
	}

	// Check Primary DB
	if err := h.dbPrimary.PingContext(ctx); err != nil {
		status["services"].(map[string]interface{})["database"].(map[string]interface{})["primary"].(map[string]string)["status"] = "error"
		status["services"].(map[string]interface{})["database"].(map[string]interface{})["primary"].(map[string]string)["message"] = err.Error()
	} else {
		status["services"].(map[string]interface{})["database"].(map[string]interface{})["primary"].(map[string]string)["status"] = "healthy"
	}

	// Check Replica DB
	if err := h.dbReplica.PingContext(ctx); err != nil {
		status["services"].(map[string]interface{})["database"].(map[string]interface{})["replica"].(map[string]string)["status"] = "error"
		status["services"].(map[string]interface{})["database"].(map[string]interface{})["replica"].(map[string]string)["message"] = err.Error()
	} else {
		status["services"].(map[string]interface{})["database"].(map[string]interface{})["replica"].(map[string]string)["status"] = "healthy"
	}

	// Check Redis
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		status["services"].(map[string]interface{})["redis"].(map[string]string)["status"] = "error"
		status["services"].(map[string]interface{})["redis"].(map[string]string)["message"] = err.Error()
	} else {
		status["services"].(map[string]interface{})["redis"].(map[string]string)["status"] = "healthy"
	}

	// Overall status
	allHealthy := true
	services := status["services"].(map[string]interface{})
	if services["database"].(map[string]interface{})["primary"].(map[string]string)["status"] != "healthy" ||
		services["database"].(map[string]interface{})["replica"].(map[string]string)["status"] != "healthy" ||
		services["redis"].(map[string]string)["status"] != "healthy" {
		allHealthy = false
	}

	status["status"] = map[string]interface{}{
		"code":    http.StatusOK,
		"healthy": allHealthy,
	}

	if !allHealthy {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
