package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	// c.JSON(http.StatusNotFound, gin.H{
	// 	"success": false,
	// 	"message": err.Error(),
	// })
	WriteFailed(c, err, http.StatusNotFound)
}

func WriteFailed(c *gin.Context, err error, errorCode int) {
	c.JSON(errorCode, gin.H{
		"success": false,
		"message": err.Error(),
	})
}
