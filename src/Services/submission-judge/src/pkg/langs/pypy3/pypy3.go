package pypy3

import (
	"io"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/utils"
)

type PyPy3 struct {
	id          string
	name        string
	compileArgs []string
	needCompile bool
	pkg.LanguageService
}

func (pypy3 PyPy3) ID() string {
	return pypy3.id
}

func (pypy3 PyPy3) DisplayName() string {
	return pypy3.name
}

func (pypy3 PyPy3) DefaultFileName() string {
	return "main.py"
}

func (pypy3 PyPy3) FileExtension() string {
	return "py"
}

func (pypy3 PyPy3) ExecutableName() string {
	return "main"
}

// Run cpp file, which is a binary file, make sure it's present in the isolate working directory
func (pypy3 PyPy3) Run(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) error {
	i.Logger.Info().Msgf("Start running source code with id: %s", req.SubmissionId)

	runArgs := pypy3.compileArgs
	runArgs = append(runArgs, []string{
		utils.GetMappedFileNamePath(pypy3.DefaultFileName()),
	}...)

	i.Logger.Info().Msgf("Start compiling source code with id: %s", req.SubmissionId)

	return pypy3.IService.Run(
		i, *rc, req, "/usr/bin/pypy3", runArgs...,
	)
}

func (pypy3 PyPy3) RunCmdStrNoStream(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) ([]string, error) {
	i.Logger.Info().Msgf("Start running source code with id: %s", req.SubmissionId)

	runArgs := pypy3.compileArgs
	runArgs = append(runArgs, []string{
		utils.GetMappedFileNamePath(pypy3.DefaultFileName()),
	}...)

	i.Logger.Info().Msgf("Start compiling source code with id: %s", req.SubmissionId)

	return pypy3.IService.RunCmdStrNoStream(
		i, *rc, req, "/usr/bin/pypy3", runArgs...,
	)
}

// Compile file, make sure the file is present inside the isolate working directory
func (pypy3 PyPy3) Compile(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error {
	// compile
	rc := domain.RunConfig{
		TimeLimit:    time.Second * 60,
		MemoryLimit:  1024 * memory.MiB,
		MaxProcesses: 200,
		InheritEnv:   true,
		Stdout:       stderr,
		Stderr:       stderr,
		Meta:         true,
	}

	runArgs := pypy3.compileArgs
	runArgs = append(runArgs, []string{
		utils.GetMappedFileNamePath(pypy3.DefaultFileName()),
	}...)

	i.Logger.Info().Msgf("Start compiling source code with id: %s", req.SubmissionId)

	return pypy3.IService.Run(
		i, rc, req, "/usr/bin/pypy3", runArgs...,
	)
}

var DefaultCompileArgs = []string{}

var pypy3 = PyPy3{
	id:          "pypy3",
	name:        "PyPy3",
	compileArgs: make([]string, 0),
}

func GetAllOptions() []pkg.Language {
	return []pkg.Language{pypy3}
}
