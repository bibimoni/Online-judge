package interactorimpl

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	checkerimpl "github.com/bibimoni/Online-judge/submission-judge/src/service/checker/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/interactor"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem/utils"
)

type InteractorServiceImpl struct{}

func NewInteractorService() interactor.InteractorService {
	return NewInteractorServiceImpl()
}

func NewInteractorServiceImpl() *InteractorServiceImpl {
	return &InteractorServiceImpl{}
}

func verifyInteractor(interactorAddr, inputAddr, outputAddr string) error {
	b, err := utils.FileExsits(interactorAddr)
	if err != nil {
		return err
	}
	if !b {
		return interactor.FileNotExits
	}

	b, err = utils.FileExsits(inputAddr)
	if err != nil {
		return err
	}
	if !b {
		return interactor.FileNotExits
	}

	b, err = utils.FileExsits(outputAddr)
	if err != nil {
		return err
	}
	if !b {
		return interactor.FileNotExits
	}

	return nil
}

func (is *InteractorServiceImpl) RunInteractor(crossRunAddr, interactorAddr, inputAddr, outputAddr, answerAddr, reportAddr string, isolateStr []string) (domain.Verdict, int, string, error) {
	log := config.GetLogger()
	log.Debug().Msgf("Run CrossRun!")
	if err := verifyInteractor(interactorAddr, inputAddr, outputAddr); err != nil {
		log.Debug().Msgf("Error when run CrossRun")
		return "", -1, "", err
	}

	interactorCmdStr := fmt.Sprintf("%s %s %s %s %s", interactorAddr, inputAddr, outputAddr, answerAddr, reportAddr)
	isolateCmdStr := fmt.Sprintf("%s", strings.Join(isolateStr, " "))
	log.Debug().Msgf("Isolate str: %s", isolateCmdStr)
	cmd := []string{"java", "-jar", crossRunAddr, isolateCmdStr, interactorCmdStr}
	log.Debug().Msgf("CMD String: %s", cmd)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	execCmd := exec.CommandContext(ctx, cmd[0], cmd[1:]...)

	combined, err := execCmd.CombinedOutput()
	msg := strings.TrimSpace(string(combined))

	// Check if timeout occurred
	if ctx.Err() == context.DeadlineExceeded {
		log.Warn().Msgf("CrossRun timed out after 15 seconds on test case")
		log.Warn().Msgf("This usually means a pipe deadlock or infinite loop")
		log.Warn().Msgf("Output so far: %s", msg)
		return domain.JUDGEMENT_FAILED, -1, "Interactive process timed out", nil
	}

	log.Debug().Msgf("Message: %s, exit code: %v", msg, err)

	if err == nil {
		return checkerimpl.MapExitCodeToVerdict(0), 0, msg, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		return checkerimpl.MapExitCodeToVerdict(code), code, msg, nil
	}

	return "", -1, "", nil
}
