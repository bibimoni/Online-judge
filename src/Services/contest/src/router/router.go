package router

import (
	"contest/src/components"
	"contest/src/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, appContext components.AppContext) {
	contestController := controller.NewContestController(appContext)

	contest := group.Group("/contest")
	contest.POST("/create", contestController.Create())
	contest.POST("/edit", contestController.Edit())
}
