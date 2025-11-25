package router

import (
	appctx "github.com/bibimoni/Online-judge/submission-judge/src/components"
	transportgetsubmission "github.com/bibimoni/Online-judge/submission-judge/src/controller/getsubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/controller/submitsubmission"
	controller_utils "github.com/bibimoni/Online-judge/submission-judge/src/controller/utils"
	"github.com/bibimoni/Online-judge/submission-judge/src/controller/websocketsubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission/interactor"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, appContext appctx.AppContext) {
	submissionInteractor, err := controller_utils.InitInteractor(appContext)
	if err != nil {
		config.GetLogger().Panic().Err(err).Msg("Can't initialize interactor for submission")
	}

	rdb := appContext.GetRedis()
	rrepo := impl.NewRedisSubmissionRepository(rdb)
	wsSubmissionInteractor := interactor.NewWSSubmissionInteractor(rrepo)

	submission := group.Group("/submission")
	submission.POST("/submit", transportsubmitsubmission.HandleSubmitSubmissionRequest(appContext, submissionInteractor))
	submission.GET("/view/:submission_id", transportgetsubmission.HandleGetSubmissionRequest(appContext, submissionInteractor))
	submission.GET("/ws", websocketsubmission.HandleSubmissionWSRequest(appContext, wsSubmissionInteractor))
	submission.GET("/problem/view/:problem_id", transportgetsubmission.HandleGetProblemSubmissionRequest(appContext, submissionInteractor))
}
