package controller

import "github.com/gin-gonic/gin"

type RequestContext struct {
	Username string
	UserId   uint64
	Role     string
	Perms    []string
}

const requestContextKey = "request_context"

func SetRequestContext(c *gin.Context, rc *RequestContext) {
	c.Set(requestContextKey, rc)
}

func GetContextRequest(c *gin.Context) (*RequestContext, bool) {
	val, exists := c.Get(requestContextKey)
	if !exists {
		return nil, false
	}
	rc, ok := val.(*RequestContext)
	return rc, ok
}
