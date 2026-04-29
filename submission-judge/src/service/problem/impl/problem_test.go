package impl_test

import (
	"testing"

	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem/impl"
)

func TestProblemServiceImpl_GetTestCaseAddr(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		problemId string
		version   string
		tcType    problem.TestCaseType
		testNum   int
		want      string
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			"test problem service",
			"445985",
			"v1",
			"INPUT",
			1,
			"_",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, err := impl.NewProblemServiceImpl()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := ps.GetTestCaseAddr(tt.problemId, tt.version, tt.tcType, tt.testNum)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetTestCaseAddr() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetTestCaseAddr() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			// if true {
			// }
			t.Errorf("GetTestCaseAddr() = %v, want %v", got, tt.want)
		})
	}
}
