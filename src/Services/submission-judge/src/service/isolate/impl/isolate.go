package impl

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/utils"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem/impl"
)

/*
isolate -b 0 --time=2 --mem=256000 --dir=/in=/problems_dir/445985/tests/input:rw
--dir=/app=/var/local/lib/isolate/0/box/submissionId:rw -i /in/01 -o /app/01 -M /app/meta --run -- /app/main
*/

type IsolateServiceImpl struct {
	problemService problem.ProblemService
}

func NewIsolateServiceImpl() (*IsolateServiceImpl, error) {
	ps, err := impl.NewProblemService()
	if err != nil {
		return nil, err
	}
	return &IsolateServiceImpl{
		problemService: ps,
	}, nil
}

func NewIsolateService() (isolateservice.IsolateService, error) {
	is, err := NewIsolateServiceImpl()
	if err != nil {
		return nil, err
	}
	return is, nil
}

type Task struct {
	SubmissionId   string
	Username       string
	Sourcecode     string
	SubmissionType string
	ProblemId      string
}

func (ir *IsolateServiceImpl) NewIsolate(id int) (*domain.Isolate, error) {
	res := &domain.Isolate{
		ID: id,
	}
	if res.Logger == nil {
		log, err := config.NewIsolateLogger(id)
		if err != nil {
			return nil, err
		}
		(*res).Logger = log
	}
	return res, nil
}

func (ir *IsolateServiceImpl) Cleanup(i *domain.Isolate) error {
	cmd := []string{"isolate", "--cg", "-b", strconv.Itoa(i.ID), "--cleanup"}
	i.Logger.Info().Msgf("Cleaning up... Running: %s", cmd)
	i.Inited = false

	return exec.Command(cmd[0], cmd[1:]...).Run()
}

func (ir *IsolateServiceImpl) Init(i *domain.Isolate) error {
	if err := ir.Cleanup(i); err != nil {
		return err
	}

	cmd := []string{"isolate", "--cg", "-b", strconv.Itoa(i.ID), "--init"}
	i.Logger.Info().Msgf("Creating isolate... Running: %s", cmd)
	i.Inited = true
	i.BoxDir = filepath.Join(isolateservice.IsolateRoot, strconv.Itoa(i.ID))

	log := config.GetLogger()
	log.Info().Msgf("Trying to tnit isolate service with id: %d", i.ID)
	return exec.Command(cmd[0], cmd[1:]...).Run()
}

// This method added input directory to the run config, this also verify if the directory exists
func (ir *IsolateServiceImpl) addInputMappedDir(rc *domain.RunConfig, problemId string) error {
	inputDirAddr, err := ir.problemService.GetTestCaseDirAddr(problemId, problem.INPUT)
	if err != nil {
		return err
	}

	rc.DirectoryMaps = append(rc.DirectoryMaps, domain.DirectoryMap{
		Inside:  "/" + isolateservice.IsolateInputDirName,
		Outside: inputDirAddr,
		Options: []domain.DirectoryMapOption{domain.AllowReadWrite},
	})
	return nil
}

// This method added working directory to the run config, so access it only via /app which makes it easier to implement logic
// This also map it's current working directory to itself
func addWorkingMappedDir(i *domain.Isolate, rc *domain.RunConfig, submissionId string) {
	workingDir := utils.GetSubmissionDir(i, submissionId)

	rc.DirectoryMaps = append(rc.DirectoryMaps, domain.DirectoryMap{
		Inside:  "/" + isolateservice.IsolateWorkingDirName,
		Outside: workingDir,
		Options: []domain.DirectoryMapOption{domain.AllowReadWrite},
	})
}

func buildArgs(i *domain.Isolate, rc domain.RunConfig, submissionId string) ([]string, error) {
	if !i.Inited {
		return []string{}, isolateservice.ErrorIsolateNotInitialized
	}
	// --cg MUST come before -b flag
	args := []string{"isolate", "--cg", "-b", strconv.Itoa(i.ID)}

	if rc.MaxProcesses > 0 {
		args = append(args, fmt.Sprintf("--processes=%d", rc.MaxProcesses))
	} else {
		args = append(args, "--processes=100")
	}

	if rc.InheritEnv {
		args = append(args, "--full-env")
	}
	for ind := range rc.Env {
		args = append(args, fmt.Sprintf("--env=%s", rc.Env[ind]))
	}
	for _, rule := range rc.DirectoryMaps {
		arg := fmt.Sprintf("--dir=%s=%s", rule.Inside, rule.Outside)
		for _, opt := range rule.Options {
			arg += ":" + string(opt)
		}
		args = append(args, arg)
	}

	if rc.TimeLimit > 0 {
		ms := rc.TimeLimit / time.Millisecond
		args = append(args, fmt.Sprintf("--time=%d.%d", ms/1000, ms%1000))
		args = append(args, fmt.Sprintf("--wall-time=%d.%d", (2*ms+1000)/1000, (2*ms+1000)%1000))
	}

	if rc.MemoryLimit > 0 {
		// args = append(args, fmt.Sprintf("--mem=%d", int(rc.MemoryLimit/memory.KiB)))
		args = append(args, fmt.Sprintf("--cg-mem=%d", int(rc.MemoryLimit/memory.KiB)))
	}

	if len(rc.Input) > 0 {
		args = append(args, "-i")
		args = append(args, rc.Input)
	}

	if len(rc.Output) > 0 {
		args = append(args, "-o")
		args = append(args, rc.Output)
	}

	if rc.Meta {
		args = append(args, "-M")
		args = append(args, utils.GetSubmissionDir(i, submissionId)+"/"+isolateservice.IsolateMetaFileName)
	}

	args = append(args, rc.Args...)
	return args, nil
}

func (ir *IsolateServiceImpl) RunBinary(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, exeName string) error {
	return ir.Run(i, rc, req, "/"+isolateservice.IsolateWorkingDirName+"/"+exeName)
}

func (ir *IsolateServiceImpl) Run(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) error {
	cmdStr, err := ir.RunCmdStrNoStream(i, rc, req, toRun, toRunArgs...)
	if err != nil {
		return err
	}
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.Stdin = rc.Stdin
	cmd.Stdout = rc.Stdout
	cmd.Stderr = rc.Stderr

	return cmd.Run()
}

func (ir *IsolateServiceImpl) RunCmdStrNoStream(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) ([]string, error) {
	if !i.Inited {
		return nil, isolateservice.ErrorIsolateNotInitialized
	}
	i.Logger.Info().Msgf("Start running command!, toRun: %s", toRun)

	addWorkingMappedDir(i, &rc, req.SubmissionId)
	ir.addInputMappedDir(&rc, req.ProblemId)
	i.Logger.Info().Msgf("Run config: %v", rc)

	args, err := buildArgs(i, rc, req.SubmissionId)
	if err != nil {
		return nil, err
	}

	args = append(args, "--run", "--", toRun)
	i.Logger.Info().Msgf("To run args: %v", toRunArgs)
	args = append(args, toRunArgs...)

	i.Logger.Info().Msgf("Running command with args: %v", args)

	return args, nil
}
