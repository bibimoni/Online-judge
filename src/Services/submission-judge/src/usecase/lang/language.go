package usecaselang

import "context"

type LanguageUsecase interface {
	GetLanguageList(ctx context.Context, _ *struct{}) []string
}

type GetLanguageListOutput struct {
	LanguageList []LanguageInfo `json:"language_list"`
}

type LanguageInfo struct {
	LangID      string `json:"id"`
	LangName    string `json:"name"`
	LangRuntime string `json:"runtime"`
}
