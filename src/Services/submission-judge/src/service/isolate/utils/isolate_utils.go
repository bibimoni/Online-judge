package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
)

// Create a new submission file with the content from @param{sourceCode} inside the isolate working directory
func CreateSubmissionSourceFile(i *domain.Isolate, sourceCode string, submissionId string, sourceCodeName string) (*os.File, error) {
	sourceCodeAddr := GetIsolateDir(i) + "/" + submissionId + "/" + sourceCodeName
	dir := filepath.Dir(sourceCodeAddr)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(sourceCodeAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create source code in: %s", sourceCodeAddr)
	}
	_, err = file.WriteString(sourceCode)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func CopyChecker(i *domain.Isolate, submissionId string, checkerLocation string) error {
	in, err := os.Open(checkerLocation)
	if err != nil {
		return err
	}
	defer in.Close()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	checkerNewAddr := GetIsolateDir(i) + "/" + submissionId + "/" + cfg.CheckerBinName

	out, err := os.Create(checkerNewAddr)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if err = out.Sync(); err != nil {
		return err
	}

	info, err := os.Stat(checkerLocation)
	if err != nil {
		return err
	}
	if err = os.Chmod(checkerNewAddr, info.Mode()); err != nil {
		return err
	}

	return nil
}

func CopyInteractor(i *domain.Isolate, submissionId string, interactorLocation string) error {
	in, err := os.Open(interactorLocation)
	if err != nil {
		return err
	}
	defer in.Close()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	interactorNewAddr := GetIsolateDir(i) + "/" + submissionId + "/" + cfg.InteractorBinName

	out, err := os.Create(interactorNewAddr)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if err = out.Sync(); err != nil {
		return err
	}

	info, err := os.Stat(interactorLocation)
	if err != nil {
		return err
	}
	if err = os.Chmod(interactorNewAddr, info.Mode()); err != nil {
		return err
	}

	return nil
}

func CopyCrossRun(i *domain.Isolate, submissionId string, crossrunLocation string) error {
	in, err := os.Open(crossrunLocation)
	if err != nil {
		return err
	}
	defer in.Close()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	crossrunNewAddr := GetIsolateDir(i) + "/" + submissionId + "/" + cfg.CrossRunJarName

	out, err := os.Create(crossrunNewAddr)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if err = out.Sync(); err != nil {
		return err
	}

	info, err := os.Stat(crossrunLocation)
	if err != nil {
		return err
	}
	if err = os.Chmod(crossrunNewAddr, info.Mode()); err != nil {
		return err
	}
	return nil
}

func GetMetaFilePath(i *domain.Isolate, submissionId string) (string, error) {
	if !i.Inited {
		return "", isolateservice.ErrorIsolateNotInitialized
	}
	return GetSubmissionDir(i, submissionId) + "/" + isolateservice.IsolateMetaFileName, nil
}

func GetIsolateDir(i *domain.Isolate) string {
	return isolateservice.IsolateRoot + strconv.Itoa(i.ID) + "/box"
}

func GetIsolateInputDir(submissionId string) string {
	return submissionId + "/" + isolateservice.IsolateInputDirName
}

func GetIsolateWorkingDir(submissionId string) string {
	return submissionId + "/" + isolateservice.IsolateWorkingDirName
}

func GetSubmissionDir(i *domain.Isolate, submissionId string) string {
	return GetIsolateDir(i) + "/" + submissionId
}

func GetMappedFileNamePath(fileName string) string {
	return "/" + isolateservice.IsolateWorkingDirName + "/" + fileName
}
