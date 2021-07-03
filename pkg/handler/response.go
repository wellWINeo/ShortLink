package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, msg string) {
	log.Printf("Error response: %s", msg)
	c.AbortWithStatusJSON(statusCode, errorResponse{msg})
}
