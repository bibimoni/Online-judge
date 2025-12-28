package pkg

import (
	"io"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
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
	SetIsolateService(isolateService isolateservice.IsolateService)
}

type LanguageService struct {
	IService isolateservice.IsolateService
}

func (langSerivce *LanguageService) SetIsolateService(isolateService isolateservice.IsolateService) {
	langSerivce.IService = isolateService
}
