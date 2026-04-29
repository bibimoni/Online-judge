package utils_test

import (
	"fmt"
	"os"
	"problem/utils"
	"testing"
)

func TestParseProblem(t *testing.T) {
	f, err := os.Open("problem.xml")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	problem, err := utils.ParseProblemStruct(123, f)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", problem)
}
