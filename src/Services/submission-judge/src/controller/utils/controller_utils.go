package controllerutils

import (
	appctx "github.com/bibimoni/Online-judge/submission-judge/src/components"
	"github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission/interactor"
)

func InitInteractor(appContext appctx.AppContext) (*interactor.SubmissionInteractor, error) {
	return interactor.NewSubmissionInteractor(
		appContext.GetSubmissionRepo(),
		appContext.GetSourcecodeRepo(),
		appContext.GetProblemService(),
		appContext.GetJudgeService(),
		appContext.GetEvalRepo(),
	), nil
}
