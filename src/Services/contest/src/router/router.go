package router

import (
	"contest/src/components"
	contestcontroller "contest/src/controller/contest"
	"contest/src/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, appContext components.AppContext) {
	contestInteractor := appContext.GetContestInteractor()
	contest := group.Group("/contest")
	{
		// auth route for contest
		auth := contest.Group("")
		auth.Use(middleware.RequireAuth())
		auth.POST("/create", contestcontroller.Create(contestInteractor))
		auth.POST("/edit", contestcontroller.Edit(contestInteractor))
		auth.PATCH("/patch/:contest_id", contestcontroller.Patch(contestInteractor))
	}
}
