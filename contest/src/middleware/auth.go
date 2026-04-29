package middleware

import (
	"contest/src/common"
	"contest/src/controller"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-Username")
		userId := c.GetHeader("X-User-Id")
		userRole := c.GetHeader("X-User-Role")
		userPerms := c.GetHeader("X-User-Permissions")
		// var perms map[string]any
		// config.GetLogger().Debug().Msgf("Auth Headers: %+v", c.Request.Header)

		if username == "" || userId == "" || userRole == "" {
			common.WriteFailedOutput(c, common.NewForbiddenError("authentication required"))
			c.Abort()
			return
		}
		convUserId, err := strconv.Atoi(userId)
		if err != nil {
			common.WriteFailedOutput(c, common.NewForbiddenError("invalid user id"))
			c.Abort()
			return
		}
		rc := controller.RequestContext{
			Username: username,
			Role:     userRole,
			Perms:    []string{userPerms},
			UserId:   uint64(convUserId),
		}
		controller.SetRequestContext(c, &rc)
		c.Next()
	}
}

func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-Username")
		userId := c.GetHeader("X-User-Id")
		userRole := c.GetHeader("X-User-Role")
		userPerms := c.GetHeader("X-User-Permissions")

		if username != "" && userId != "" && userRole != "" {
			convUserId, err := strconv.Atoi(userId)
			if err == nil {
				rc := controller.RequestContext{
					Username: username,
					Role:     userRole,
					Perms:    []string{userPerms},
					UserId:   uint64(convUserId),
				}
				controller.SetRequestContext(c, &rc)
			}
		}
		c.Next()
	}
}
