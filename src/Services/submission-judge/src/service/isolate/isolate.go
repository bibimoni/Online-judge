package isolateservice

import (
	"errors"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
)

// IsolateRoot is the root directory structure isolate is using.
var IsolateRoot = "/var/local/lib/isolate/"
var IsolateInputDirName = "in"
var IsolateWorkingDirName = "app"
var IsolateMetaFileName = "meta"

var ErrorIsolateNotInitialized = errors.New("initialize the isolate first")

type IsolateService interface {
	NewIsolate(id int) (*domain.Isolate, error)
	Cleanup(i *domain.Isolate) error
	Init(i *domain.Isolate) error
	Run(i *domain.Isolate, rc domain.RunConfig, req *SubmissionRequest, toRun string, toRunArgs ...string) error
	RunBinary(i *domain.Isolate, rc domain.RunConfig, req *SubmissionRequest, exeName string) error
	RunCmdStrNoStream(i *domain.Isolate, rc domain.RunConfig, req *SubmissionRequest, toRun string, toRunArgs ...string) ([]string, error)
}

type SubmissionRequest struct {
	SubmissionId   string
	Username       string
	Sourcecode     string
	SubmissionType domain.SubmissionType
	ProblemId      string
	LanguageId     string
	EvalId         string
}
