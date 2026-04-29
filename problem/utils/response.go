package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func WriteSuccessOutput[T any](c *fiber.Ctx, output *T, err error) {
	if err != nil {
		WriteFailedOutput(c, err)
		return
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"success": true,
	// 	"data":    output,
	// })

	c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    output,
	})
}

func WriteCreatedOutput[T any](c *fiber.Ctx, output *T, err error) {
	if err != nil {
		WriteFailedOutput(c, err)
		return
	}
	c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    output,
	})
}

func WriteFailedOutput(c *fiber.Ctx, err error) {
	// c.JSON(http.StatusNotFound, gin.H{
	// 	"success": false,
	// 	"message": err.Error(),
	// })
	WriteFailed(c, err, http.StatusNotFound)
}

func WriteFailed(c *fiber.Ctx, err error, errorCode int) {
	c.Status(errorCode).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}
