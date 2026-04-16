package domain

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type EvaluationResult struct {
	Id              bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	SubmissionId    bson.ObjectID   `json:"submission_id,omitempty" bson:"submission_id"`
	Verdict         Verdict         `json:"verdict,omitempty" bson:"verdict,omitempty"`
	VerdictCase     []Verdict       `json:"verdict_case,omitempty" bson:"verdict_case,omitempty"`
	CpuTime         float64         `json:"cpu_time,omitempty" bson:"cpu_time,omitempty"`
	CpuTimeCase     []float64       `json:"cpu_time_case,omitempty" bson:"cpu_time_case,omitempty"`
	MemoryUsage     memory.Memory   `json:"memory_usage,omitempty" bson:"memory_usage,omitempty"`
	MemoryUsageCase []memory.Memory `json:"memory_usage_case,omitempty" bson:"memory_usage_case,omitempty"`
	NSuccess        int             `json:"n_success,omitempty" bson:"n_success,omitempty"`
	Outputs         []string        `json:"outputs,omitempty" bson:"outputs,omitempty"`

	Score        float64   `json:"score,omitempty" bson:"score,omitempty"`
	MaxScore     float64   `json:"max_score,omitempty" bson:"max_score,omitempty"`
	ScoreCase    []float64 `json:"score_case,omitempty" bson:"score_case,omitempty"`
	MaxScoreCase []float64 `json:"max_score_case,omitempty" bson:"max_score_case,omitempty"`

	Points     int   `json:"points,omitempty" bson:"points,omitempty"`
	PointsCase []int `json:"points_case,omitempty" bson:"points_case,omitempty"`

	TimestampFinish int              `json:"timestamp_finish,omitempty" bson:"timestamp_finish,omitempty"`
	Message         string           `json:"message,omitempty" bson:"message,omitempty"`
	EvalStatus      SubmissionStatus `json:"eval_status,omitempty" bson:"eval_status"`
	NCases          int              `json:"n_cases,omitempty" bson:"n_cases"`
	TL              int              `json:"tl,omitempty" bson:"tl"`
	ML              memory.Memory    `json:"ml,omitempty" bson:"ml"`
}

type TestGroup struct {
	Name         string      `json:"name,omitempty" bson:"name,omitempty"`
	Scoring      ScoringType `json:"scoring,omitempty" bson:"scoring,omitempty"`
	Score        float64     `json:"score,omitempty" bson:"score,omitempty"`
	MaxScore     float64     `json:"max_score,omitempty" bson:"max_score,omitempty"`
	TestIndices  []int       `json:"test_indices,omitempty" bson:"test_indices,omitempty"`
	Dependencies []string    `json:"dependencies,omitempty" bson:"dependencies,omitempty"`
}

type ScoringType string

const (
	ScoringSum   ScoringType = "SUM"   // Sum all test scores
	ScoringGroup ScoringType = "GROUP" // All-or-nothing (complete-group)
	ScoringMin   ScoringType = "MIN"   // Minimum score
)

type Verdict string

const (
	ACCEPTED              Verdict = "ACCEPTED"
	COMPILATION_ERROR     Verdict = "COMPILATION_ERROR"
	REJECTED              Verdict = "REJECTED"
	RUNTIME_ERROR         Verdict = "RUNTIME_ERROR"
	TIME_LIMIT_EXCEEDED   Verdict = "TIME_LIMIT_EXCEEDED"
	MEMORY_LIMIT_EXCEEDED Verdict = "MEMORY_LIMIT_EXCEEDED"
	WRONG_ANSWER          Verdict = "WRONG_ANSWER"
	JUDGEMENT_FAILED      Verdict = "JUDGEMENT_FAILED"
	PRESENTATION_ERROR    Verdict = "PRESENTATION_ERROR"
	FAIL                  Verdict = "FAIL"
	POINTS                Verdict = "POINTS"
	PARTIAL_RESULT        Verdict = "PARTIAL_RESULT"
	UNEXPECTED_EOF        Verdict = "UNEXPECTED_EOF"
	DIRT                  Verdict = "DIRT"
)

type SubmissionStatus string

const (
	PENDING  SubmissionStatus = "PENDING"
	JUDGING  SubmissionStatus = "JUDGING"
	FINISHED SubmissionStatus = "FINISHED"
)
