package pkg

import (
	"io"
	"strings"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
)

type Language interface {
	ID() string
	DisplayName() string
	DefaultFileName() string
	ExecutableName() string
	FileExtension() string
	Run(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) error
	RunCmdStrNoStream(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) ([]string, error)
	Compile(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error
	InitLanguageService(isolateService isolateservice.IsolateService, language Language)
	GetCompilerBin() string
	GetCompileArgs() []string
	GetLanguageRuntime() string
}

type LanguageService struct {
	impl     Language
	IService isolateservice.IsolateService
}

func (langSerivce *LanguageService) InitLanguageService(isolateService isolateservice.IsolateService, language Language) {
	langSerivce.IService = isolateService
	langSerivce.impl = language
}

func (langSerivce *LanguageService) GetLanguageRuntime() string {
	config.GetLogger().Debug().Msgf("compile args: %v", langSerivce.impl.GetCompileArgs())
	runtime := append([]string{langSerivce.impl.GetCompilerBin()}, langSerivce.impl.GetCompileArgs()...)
	return strings.Join(runtime, " ")
}
