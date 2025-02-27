package domain

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) APIResponse {
	return APIResponse{Status: "success", Message: "OK", Data: data}
}

func Error(message string) APIResponse {
	return APIResponse{Status: "error", Message: message}
}

func JSONResponse(c *gin.Context, httpCode int, apiResponse APIResponse) {
	c.JSON(httpCode, apiResponse)
}
