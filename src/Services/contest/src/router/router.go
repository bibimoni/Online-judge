package router

import (
	"contest/src/components"
	contestcontroller "contest/src/controller/contest"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, appContext components.AppContext) {
	contestInteractor := appContext.GetContestInteractor()
	contest := group.Group("/contest")
	contest.POST("/create", contestcontroller.Create(contestInteractor))
	contest.POST("/edit", contestcontroller.Edit(contestInteractor))
	contest.PATCH("/patch/:contest_id", contestcontroller.Patch(contestInteractor))
}
