package python3

import (
	"io"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/utils"
)

type Python3 struct {
	id          string
	name        string
	compileArgs []string
	needCompile bool
	pkg.LanguageService
}

func (python Python3) ID() string {
	return python.id
}

func (python Python3) DisplayName() string {
	return python.name
}

func (python Python3) DefaultFileName() string {
	return "main.py"
}

func (python Python3) FileExtension() string {
	return "py"
}

func (python Python3) ExecutableName() string {
	return "main"
}

func (python Python3) GetCompileArgs() []string {
	return python.compileArgs
}

func (python Python3) GetCompilerBin() string {
	return "/usr/bin/python3"
}

// Run cpp file, which is a binary file, make sure it's present in the isolate working directory
func (python Python3) Run(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) error {
	i.Logger.Info().Msgf("Start running source code with id: %s", req.SubmissionId)

	runArgs := python.compileArgs
	runArgs = append(runArgs, []string{
		utils.GetMappedFileNamePath(python.DefaultFileName()),
	}...)

	i.Logger.Info().Msgf("Start compiling source code with id: %s", req.SubmissionId)

	return python.IService.Run(
		i, *rc, req, python.GetCompilerBin(), runArgs...,
	)
}

func (python Python3) RunCmdStrNoStream(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) ([]string, error) {
	i.Logger.Info().Msgf("Start running source code with id: %s", req.SubmissionId)

	runArgs := python.compileArgs
	runArgs = append(runArgs, []string{
		utils.GetMappedFileNamePath(python.DefaultFileName()),
	}...)

	i.Logger.Info().Msgf("Start compiling source code with id: %s", req.SubmissionId)

	return python.IService.RunCmdStrNoStream(
		i, *rc, req, python.GetCompilerBin(), runArgs...,
	)
}

// Compile file, make sure the file is present inside the isolate working directory
func (python Python3) Compile(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error {
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

	runArgs := python.compileArgs
	runArgs = append(runArgs, []string{
		utils.GetMappedFileNamePath(python.DefaultFileName()),
	}...)

	i.Logger.Info().Msgf("Start compiling source code with id: %s", req.SubmissionId)

	return python.IService.Run(
		i, rc, req, python.GetCompilerBin(), runArgs...,
	)
}

var DefaultCompileArgs = []string{}

var python3 = Python3{
	id:          "python3",
	name:        "Python3",
	compileArgs: make([]string, 0),
}

func GetAllOptions() []pkg.Language {
	return []pkg.Language{&python3}
}
