package handler

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func JSONSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func JSONError(c *gin.Context, status int, err error) {
	c.JSON(status, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func JSONErrorMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, Response{
		Success: false,
		Error:   msg,
	})
}
