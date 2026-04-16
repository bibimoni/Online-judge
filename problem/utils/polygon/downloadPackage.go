package polygon

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"problem/models"
	"problem/utils"
	"strconv"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/xyproto/unzip"
)

func checkVersionIntegrity(problemId uint64) error {
	versions := make([]bool, 101)
	for versionNum := 1; versionNum <= 100; versionNum++ {
		destpath := fmt.Sprintf("%s/%d/v%d", os.Getenv("PROBLEM_STORAGE_DIR"), problemId, versionNum)

		_, err := os.Stat(destpath)
		if os.IsNotExist(err) {
			versions[versionNum] = false
		} else {
			versions[versionNum] = true
		}

		// if err != nil {
		// 	fmt.Printf("An error occurred: %v\n", err)
		// 	return
		// }

		if versionNum > 1 && (versions[versionNum] == true && versions[versionNum-1] == false) {
			return errors.New("problem versions' integrity check failed")
		}
	}

	return nil
}

func getNextVersionNumber(problemId uint64) (uint64, error) {
	problemDir := fmt.Sprintf("%s/%d", os.Getenv("PROBLEM_STORAGE_DIR"), problemId)
	if _, err := os.Stat(problemDir); os.IsNotExist(err) {
		if err := os.MkdirAll(problemDir, os.ModePerm); err != nil {
			return 0, err
		}
	}

	var targetVersion uint64 = 0

	for versionNum := uint64(1); versionNum <= 100; versionNum++ {
		destpath := fmt.Sprintf("%s/%d/v%d", os.Getenv("PROBLEM_STORAGE_DIR"), problemId, versionNum)

		_, err := os.Stat(destpath)
		if os.IsNotExist(err) {
			targetVersion = versionNum
			break
		}
	}

	if targetVersion == 0 {
		return 0, errors.New("maximum number of versions reached")
	}

	return targetVersion, nil
}

func GetLatestVersionNumber(problemId uint64) (uint64, error) {
	var latestVersion uint64 = 0

	for versionNum := 1; versionNum <= 100; versionNum++ {
		destpath := fmt.Sprintf("%s/%d/v%d", os.Getenv("PROBLEM_STORAGE_DIR"), problemId, versionNum)

		_, err := os.Stat(destpath)
		if os.IsNotExist(err) {
			break
		}
		latestVersion = uint64(versionNum)
	}

	if latestVersion == 0 {
		return 0, errors.New("no version found")
	}

	return latestVersion, nil
}

/*
DownloadPackge() - Download Polygon problems
- Problems are stored at $STORAGE_DIR/$problemId/ with structure following the README.md
- It removes all downloaded packages of the given problem before downloading the specified package

FUTURE:
*/
func DownloadPackage(problemId uint64, packageId uint64) error {
	if err := checkVersionIntegrity(problemId); err != nil {
		return err
	}

	params := map[string]string{
		"problemId": strconv.Itoa(int(problemId)),
		"packageId": strconv.Itoa(int(packageId)),
		"apiKey":    os.Getenv("POLYGON_API_KEY"),
		"type":      "linux",
		"time":      fmt.Sprintf("%d", time.Now().Unix()),
	}
	resp, err := polygonApiCall("problem.package", params)
	if err != nil {
		return fmt.Errorf("error making Polygon api call: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error parsing response body: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}

	// Extract to a temporary directory
	f, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return fmt.Errorf("error creating temporary zip: %s", err.Error())
	}
	defer f.Close()

	if _, err = f.Write(body); err != nil {
		return fmt.Errorf("error getting response: %s", err.Error())
	}

	// Handle versioning
	nextVersion, err := getNextVersionNumber(problemId)
	if err != nil {
		return err
	}

	destpath := fmt.Sprintf("%s/%s/v%d", os.Getenv("PROBLEM_STORAGE_DIR"), params["problemId"], nextVersion)
	if err := os.RemoveAll(destpath); err != nil {
		return err
	}
	if err := os.Mkdir(destpath, os.ModePerm); err != nil {
		return fmt.Errorf("error making temporary directory: %s", err.Error())
	}

	tempdir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tempdir)

	if err = unzip.Extract(f.Name(), tempdir); err != nil {
		return err
	}

	var xml *os.File
	if xml, err = os.Open(tempdir + "/problem.xml"); err != nil {
		return err
	}
	defer xml.Close()

	var problem models.Problem
	if problem, err = utils.ParseProblemStruct(problemId, xml); err != nil {
		return err
	}

	var errBuffer bytes.Buffer

	cmd := exec.Command("scripts/gen_statement/main.sh", tempdir, destpath)
	cmd.Stderr = &errBuffer
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating statement: %s", errBuffer.String())
	}

	// compile interactor (in interactive problems)
	f, err = os.Open(tempdir + "/problem.xml")
	if err != nil {
		return err
	}
	defer f.Close()

	doc, err := xmlquery.Parse(f)
	if err != nil {
		return fmt.Errorf("error parsing problem.xml: %s", err.Error())
	}

	checker_file := tempdir + "/" + xmlquery.FindOne(doc, "/problem/assets/checker/source").SelectAttr("path")
	cmd = exec.Command("scripts/compile_checker/main.sh", tempdir, checker_file, destpath)
	cmd.Stderr = &errBuffer
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error compiling checker: %s", errBuffer.String())
	}

	interactor_file := xmlquery.FindOne(doc, "/problem/assets/interactor")
	if interactor_file != nil {
		cmd = exec.Command("scripts/handle_interactive_problem/compile_interactor.sh", tempdir, destpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error compiling interactor: %s", errBuffer.String())
		}

		cmd = exec.Command("scripts/handle_interactive_problem/get_files.sh", tempdir, destpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error getting interactive-related files: %s", errBuffer.String())
		}

		cmd = exec.Command("scripts/handle_interactive_problem/get_tests.sh", tempdir, destpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error getting interactive tests: %s", errBuffer.String())
		}

		problem.IsInteractive = true
	} else {
		cmd := exec.Command("scripts/get_tests/main.sh", tempdir, destpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error getting tests: %s", errBuffer.String())
		}
	}

	if err := utils.SaveProblemToJson(problem, destpath+"/problem.json"); err != nil {
		return fmt.Errorf("error saving to problem.json: %s", err.Error())
	}

	return nil
}
