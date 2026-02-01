package problemservice

import "context"

// ProblemService is the same service as the problem service in submission-judge service
type ProblemService interface {
	Get(ctx context.Context, id string) (*ProblemServiceGetOutput, error)
	GetProblemMaxScores(problem *ProblemServiceGetOutput) float64
}

type ProblemServiceGetOutput struct {
	Id            string   `json:"ID,omitempty"`
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
