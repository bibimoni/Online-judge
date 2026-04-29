package utils_test

import (
	"log"
	"os"
	"problem/utils"
	"problem/utils/polygon"
	"testing"

	"github.com/joho/godotenv"
)

func TestSuccessfulDownload(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	result := polygon.DownloadPackage(332909, 1154548)
	var expected error = nil

	if result != expected {
		t.Errorf("TestSuccessfulDownload expected %s; got %s", expected, result)
	}
}

func TestParseProblemStruct_OI(t *testing.T) {
	problemId := uint64(332909)
	var xml *os.File
	var err error
	if xml, err = os.Open("332909.xml"); err != nil {
		t.Error(err)
	}

	if problem, err := utils.ParseProblemStruct(problemId, xml); err != nil {
		t.Error(err)
	} else {
		if problem.ScoringMode != "OI" {
			t.Errorf("wrong ScoringType")
		}
	}
}

func TestParseProblemStruct_ICPC(t *testing.T) {
	problemId := uint64(466874)
	var xml *os.File
	var err error
	if xml, err = os.Open("466874.xml"); err != nil {
		t.Error(err)
	}

	if problem, err := utils.ParseProblemStruct(problemId, xml); err != nil {
		t.Error(err)
	} else {
		if problem.ScoringMode != "ICPC" {
			t.Errorf("wrong ScoringType")
		}
	}
}
