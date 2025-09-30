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

/*
DownloadPackge() - Download Polygon problems
- Problems are stored at $STORAGE_DIR/$problemId/ with structure following the README.md
- It removes all downloaded packages of the given problem before downloading the specified package

FUTURE:
*/
func DownloadPackage(problemId uint64, packageId uint64) error {
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

	dirpath := fmt.Sprintf("%s/%s", os.Getenv("PROBLEM_STORAGE_DIR"), params["problemId"])
	if err := os.RemoveAll(dirpath); err != nil {
		return err
	}
	if err := os.Mkdir(dirpath, os.ModePerm); err != nil {
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

	cmd := exec.Command("scripts/gen_statement/main.sh", tempdir, dirpath)
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
	cmd = exec.Command("scripts/compile_checker/main.sh", tempdir, checker_file, dirpath)
	cmd.Stderr = &errBuffer
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error compiling checker: %s", errBuffer.String())
	}

	interactor_file := xmlquery.FindOne(doc, "/problem/assets/interactor")
	if interactor_file != nil {
		cmd = exec.Command("scripts/handle_interactive_problem/compile_interactor.sh", tempdir, dirpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error compiling interactor: %s", errBuffer.String())
		}

		cmd = exec.Command("scripts/handle_interactive_problem/get_files.sh", tempdir, dirpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error getting interactive-related files: %s", errBuffer.String())
		}

		cmd = exec.Command("scripts/handle_interactive_problem/get_tests.sh", tempdir, dirpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error getting interactive tests: %s", errBuffer.String())
		}

		problem.IsInteractive = true
	} else {
		cmd := exec.Command("scripts/get_tests/main.sh", tempdir, dirpath)
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error getting tests: %s", errBuffer.String())
		}
	}

	if err := utils.SaveProblemToJson(problem, dirpath+"/problem.json"); err != nil {
		return fmt.Errorf("error saving to problem.json: %s", err.Error())
	}

	return nil
}
