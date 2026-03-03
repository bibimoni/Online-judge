package problem

import (
	"context"
)

// This will handle the calling process, the retrieving process of getting problem
// served by the Problem service (the name might be a little confusing)
type ProblemService interface {
	Get(ctx context.Context, id string) (*ProblemServiceGetOutput, error)
	GetTestCaseAddr(problemId string, tcType TestCaseType, testNum int) (string, error)
	GetTestCaseDirAddr(problemId string, tcType TestCaseType) (string, error)
	GetCheckerAddr(problemId string) (string, error)
	GetInteractorAddr(problemId string) (string, error)
	GetCrossRunAddr(problemId string) (string, error)
}

type TestCaseType string

const (
	INPUT  TestCaseType = "INPUT"
	OUTPUT TestCaseType = "OUTPUT"
)

type ProblemServiceGetOutput struct {
	ID            string   `json:"ID,omitempty"`
	ProblemId     int64    `json:"problem-id,omitempty"`
	Name          string   `json:"name,omitempty"`
	ShortName     string   `json:"short-name,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	TestNum       int      `json:"test-num,omitempty"`
	TimeLimit     int      `json:"time-limit,omitempty"`
	MemoryLimit   int64    `json:"memory-limit,omitempty"`
	IsInteractive bool     `json:"is-interactive,omitempty"`

	ScoringMode   string         `json:"scoring-mode,omitempty"`
	TestGroups    []ProblemGroup `json:"test-groups,omitempty"`
	TestMaxScores []float64      `json:"test-max-scores,omitempty"`
}

type ProblemGroup struct {
	Name         string   `json:"name,omitempty"`
	Scoring      string   `json:"scoring,omitempty"`
	MaxScore     float64  `json:"max-score,omitempty"`
	TestIndices  []int    `json:"test-indices,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}
