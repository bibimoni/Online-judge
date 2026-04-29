package usecase

import (
	"context"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"github.com/gorilla/websocket"
)

type WSSubmissionUsecase interface {
	SubmissionStatus(ctx context.Context, input *WSSubmissionInput, out chan<- *WSSubmissionResponse) (*WSSubmissionResponse, error)
}

type (
	WSSubmissionInput struct {
		ProblemId    string
		Username     string
		SubmissionId string
		Ws           *websocket.Conn
	}

	WSSubmissionResponse struct {
		Username        string                  `json:"username,omitempty"`
		SubmissionId    string                  `json:"submission_id,omitempty"`
		ProblemId       string                  `json:"problem_id,omitempty"`
		Timestamp       time.Time               `json:"timestamp"`
		Language        string                  `json:"language,omitempty"`
		Verdict         domain.Verdict          `json:"verdict,omitempty"`
		VerdictCase     []domain.Verdict        `json:"verdict_case,omitempty"`
		CpuTime         float64                 `json:"cpu_time,omitempty"`
		CpuTimeCase     []float64               `json:"cpu_time_case,omitempty"`
		MemoryUsage     memory.Memory           `json:"memory_usage,omitempty"`
		MemoryUsageCase []memory.Memory         `json:"memory_usage_case,omitempty"`
		NSuccess        int                     `json:"n_success,omitempty"`
		Outputs         []string                `json:"outputs,omitempty"`
		Points          int                     `json:"points,omitempty"`
		PointsCase      []int                   `json:"points_case,omitempty"`
		Message         string                  `json:"message,omitempty"`
		EvalStatus      domain.SubmissionStatus `json:"eval_status,omitempty"`
	}
)
