package router

import (
	appctx "github.com/bibimoni/Online-judge/submission-judge/src/components"
	transportgetlanguages "github.com/bibimoni/Online-judge/submission-judge/src/controller/getlanguages"
	transportgetsubmission "github.com/bibimoni/Online-judge/submission-judge/src/controller/getsubmission"
	transporthealth "github.com/bibimoni/Online-judge/submission-judge/src/controller/health"
	"github.com/bibimoni/Online-judge/submission-judge/src/controller/submitsubmission"
	controllerutils "github.com/bibimoni/Online-judge/submission-judge/src/controller/utils"
	"github.com/bibimoni/Online-judge/submission-judge/src/controller/websocketsubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	interactorlang "github.com/bibimoni/Online-judge/submission-judge/src/usecase/lang/interactor"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission/interactor"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, appContext appctx.AppContext) {
	submissionInteractor, err := controllerutils.InitInteractor(appContext)
	if err != nil {
		config.GetLogger().Panic().Err(err).Msg("Can't initialize interactor for submission")
	}

	wsSubmissionInteractor := interactor.NewWSSubmissionInteractor(appContext.GetRedisRepo())
	languageInteractor := interactorlang.NewLanguageInteractor(store.DefaultStore)

	submission := group.Group("/submission")
	submission.POST("/submit", transportsubmitsubmission.HandleSubmitSubmissionRequest(submissionInteractor))
	submission.GET("/view/:submission_id", transportgetsubmission.HandleGetSubmissionRequest(submissionInteractor))
	submission.GET("/ws", websocketsubmission.HandleSubmissionWSRequest(wsSubmissionInteractor))
	submission.GET("/problem/view/:problem_id", transportgetsubmission.HandleGetProblemSubmissionRequest(submissionInteractor))
	submission.GET("/lang/all", transportgetlanguages.HandleGetLanguageListRequest(languageInteractor))
	submission.POST("/internal/rejudge", transportsubmitsubmission.HandleRejudgeSubmissionRequest(submissionInteractor))
	submission.GET("/health", transporthealth.HandleHealth())
}
