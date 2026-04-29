package interactorlang

import (
	"context"

	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	usecaselang "github.com/bibimoni/Online-judge/submission-judge/src/usecase/lang"
)

type LanguageInteractor struct {
	sService store.StoreService
}

func NewLanguageInteractor(
	store store.StoreService,
) *LanguageInteractor {
	return &LanguageInteractor{
		sService: store,
	}
}

func (li *LanguageInteractor) GetLanguageList(ctx context.Context, input *struct{}) (*usecaselang.GetLanguageListOutput, error) {
	langList := li.sService.List()
	langInfos := make([]usecaselang.LanguageInfo, 0, len(langList))
	for _, lang := range langList {
		langInfo := usecaselang.LanguageInfo{
			LangID:      lang.ID(),
			LangName:    lang.DisplayName(),
			LangRuntime: lang.GetLanguageRuntime(),
		}
		langInfos = append(langInfos, langInfo)
	}
	return &usecaselang.GetLanguageListOutput{
		LanguageList: langInfos,
	}, nil
}
