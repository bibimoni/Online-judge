package interactor

import (
	"errors"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
)

type InteractorService interface {
	// This function required jre in system and will run the CrossRun.jar file
	RunInteractor(crossRunAddr, interactorAddr, inputAddr, outputAddr, answerAddr, reportAddr string, isolateStr []string) (domain.Verdict, int, string, error)
}

var FileNotExits = errors.New("File not exists on the system")
