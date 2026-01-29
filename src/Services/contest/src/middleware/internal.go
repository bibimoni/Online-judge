package middleware

import (
	"contest/src/common"
	"contest/src/infrastructure/config"

	"github.com/gin-gonic/gin"
)

func Internal() gin.HandlerFunc {
	return func(c *gin.Context) {
		internalSecret := c.GetHeader("X-Internal-Secret")
		cfg, err := config.Load()
		if err != nil {
			common.WriteFailedOutput(c, common.NewInternalServerError("can't load config"))
			c.Abort()
			return
		}

		if internalSecret != cfg.InternalSecret {
			common.WriteFailedOutput(c, common.NewForbiddenError("forbidden"))
			c.Abort()
			return
		}

		c.Next()
	}
}
