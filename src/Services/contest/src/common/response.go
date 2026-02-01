package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusOK struct {
	Status string `json:"status"`
}

func WriteSuccessOutput[T any](c *gin.Context, output *T, err error) {
	if err != nil {
		WriteFailedOutput(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

func WriteCreatedOutput[T any](c *gin.Context, output *T, err error) {
	if err != nil {
		WriteFailedOutput(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    output,
	})
}

func WriteFailedOutput(c *gin.Context, err error) {
	var appError *AppError
	if errors.As(err, &appError) {
		WriteFailed(c, appError, appError.StatusCode)
		return
	}
	WriteFailed(c, err, http.StatusBadRequest)
}

func WriteFailed(c *gin.Context, err error, errorCode int) {
	c.JSON(errorCode, gin.H{
		"success": false,
		"message": err.Error(),
	})
}
