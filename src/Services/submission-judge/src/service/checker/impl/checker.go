package checkerimpl

import (
	"errors"
	"fmt"
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
func (cs *CheckerServiceImpl) RunChecker(checkerAddr, inputAddr, outputAddr, answerAddr string) (domain.Verdict, int, string, float64, error) {
	log := config.GetLogger()
	log.Debug().Msgf("Run checker with these files: checker - %s, input - %s, output - %s, ans - %s", checkerAddr, inputAddr, outputAddr, answerAddr)
	if err := verifyChecker(checkerAddr, inputAddr, outputAddr, answerAddr); err != nil {
		log.Debug().Msgf("Error when run checker: %v", err)
		return "", -1, "", 0.0, err
	}

	cmd := []string{checkerAddr, inputAddr, outputAddr, answerAddr}
	combined, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	msg := strings.TrimSpace(string(combined))
	log.Debug().Msgf("Message: %s, exit code: %v", msg, err)

	var score float64 = 0.0

	if err == nil {
		return MapExitCodeToVerdict(0), 0, msg, 100.0, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		verdict := MapExitCodeToVerdict(code)

		if code == 5 || code == 7 {
			fmt.Sscanf(msg, "%f", &score)
		} else if code >= 16 && code <= 116 {
			score = float64(code - 16)
		} else if verdict == domain.ACCEPTED {
			score = 100.0
		}
		return verdict, code, msg, score, nil
	}
	return "", -1, "", 0.0, err
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
