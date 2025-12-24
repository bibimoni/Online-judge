package checker

import (
	"errors"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
)

type CheckerService interface {
	RunChecker(checkerAddr string, inputAddr string, outputAddr string, answerAddr string) (domain.Verdict, int, string, float64, error)
}

var FileNotExits = errors.New("File not exists on the system")
