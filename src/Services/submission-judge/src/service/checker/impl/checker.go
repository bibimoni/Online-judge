package checkerimpl

import (
	"errors"
	"os/exec"
	"strings"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/checker"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem/utils"
)

type CheckerServiceImpl struct {
}

func NewCheckerService() checker.CheckerService {
	return NewCheckerServiceImpl()
}

func NewCheckerServiceImpl() *CheckerServiceImpl {
	return &CheckerServiceImpl{}
}

func verifyChecker(checkerAddr, inputAddr, outputAddr, answerAddr string) error {
	b, err := utils.FileExsits(checkerAddr)
	if err != nil {
		return err
	}
	if !b {
		return checker.FileNotExits
	}

	b, err = utils.FileExsits(inputAddr)
	if err != nil {
		return err
	}
	if !b {
		return checker.FileNotExits
	}

	b, err = utils.FileExsits(outputAddr)
	if err != nil {
		return err
	}
	if !b {
		return checker.FileNotExits
	}

	b, err = utils.FileExsits(answerAddr)
	if err != nil {
		return err
	}
	if !b {
		return checker.FileNotExits
	}

	return nil
}
func (cs *CheckerServiceImpl) RunChecker(checkerAddr, inputAddr, outputAddr, answerAddr string) (domain.Verdict, int, string, error) {
	log := config.GetLogger()
	log.Debug().Msgf("Run checker with these files: checker - %s, input - %s, output - %s, ans - %s", checkerAddr, inputAddr, outputAddr, answerAddr)
	if err := verifyChecker(checkerAddr, inputAddr, outputAddr, answerAddr); err != nil {
		log.Debug().Msgf("Error when run checker: %v", err)
		return "", -1, "", err
	}

	cmd := []string{checkerAddr, inputAddr, outputAddr, answerAddr}
	combined, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	msg := strings.TrimSpace(string(combined))
	log.Debug().Msgf("Message: %s, exit code: %v", msg, err)

	if err == nil {
		return MapExitCodeToVerdict(0), 0, msg, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		return MapExitCodeToVerdict(code), code, msg, nil
	}
	return "", -1, "", err
}

func MapExitCodeToVerdict(code int) domain.Verdict {
	switch code {
	case 0:
		return domain.ACCEPTED
	case 1:
		return domain.WRONG_ANSWER
	case 2:
		return domain.PRESENTATION_ERROR
	case 3:
		return domain.FAIL
	case 4:
		return domain.DIRT
	case 5, 7:
		return domain.POINTS
	case 8:
		return domain.UNEXPECTED_EOF
	default:
		if code >= 16 {
			return domain.PARTIAL_RESULT
		}
		return domain.JUDGEMENT_FAILED
	}
}
