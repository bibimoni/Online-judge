package transportgetlanguages

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	helper "github.com/bibimoni/Online-judge/submission-judge/src/controller"
	interactorlang "github.com/bibimoni/Online-judge/submission-judge/src/usecase/lang/interactor"
	"github.com/gin-gonic/gin"
)

func HandleGetLanguageListRequest(languageInteractor *interactorlang.LanguageInteractor) gin.HandlerFunc {
	return common.InvokeUseCase(
		toGetLanguageType,
		languageInteractor.GetLanguageList,
		helper.WriteSuccessOutput,
	)
}

// this usecase don't need input from request
func toGetLanguageType(c *gin.Context) (input *struct{}, err error) { return }
