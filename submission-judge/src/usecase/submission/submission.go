package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
)

type SubmissionUsecase interface {
	SubmitSubmission(ctx context.Context, input *SubmitSubmissionInput) (output *SubmitSubmissionResponse, err error)
	GetSubmission(ctx context.Context, input *GetSubmissionInput) (output *GetSubmissionOutput, err error)
	GetProblemSubmission(ctx context.Context, input *GetProblemSubmissionInput)
	RejudgeSubmission(ctx context.Context, input *RejudgeSubmissionInput) error
	InternalContestSubmitSubmission(ctx context.Context, input *InternalContestSubmitSubmissionInput) (output *SubmitSubmissionResponse, err error)
}

var ErrProblemLocked = errors.New("problem is locked in an active contest")

type (
	InternalContestSubmitSubmissionInput struct {
		ProblemId      string `json:"problem_id,omitempty"`
		Code           string `json:"code,omitempty"`
		Username       string
		ContestId      string                `json:"contest_id,omitempty"`
		LanguageId     string                `json:"language,omitempty"`
		SubmitAt       time.Time             `json:"submit_at"`
		SubmissionType domain.SubmissionType `json:"submission_type,omitempty"`
	}

	SubmitSubmissionInput struct {
		Role           common.RoleName
		Username       string
		ProblemId      string                `json:"problem_id,omitempty"`
		Code           string                `json:"code,omitempty"`
		LanguageId     string                `json:"language,omitempty"`
		SubmissionType domain.SubmissionType `json:"submission_type,omitempty"`
	}

	SubmitSubmissionResponse struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	}

	GetSubmissionInput struct {
		SubmissionId string
	}

	GetSubmissionOutput struct {
		ProblemId       string                  `json:"problem_id,omitempty"`
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
		NCases          int                     `json:"n_cases,omitempty"`
		TL              int                     `json:"tl,omitempty"`
		ML              memory.Memory           `json:"ml,omitempty"`
		Username        string                  `json:"username,omitempty"`
		Timestamp       time.Time               `json:"timestamp"`
		Type            domain.SubmissionType   `json:"type,omitempty"`
		Language        string                  `json:"language,omitempty"`
		SourceCode      string                  `json:"source_code,omitempty"`
		EvalStatus      domain.SubmissionStatus `json:"eval_status,omitempty"`
	}

	GetProblemSubmissionInput struct {
		ProblemId string
	}

	GetProblemSubmissionOutput struct {
		Submissions []usecase.WSSubmissionResponse
	}

	RejudgeSubmissionInput struct {
		SubmissionIds []string `json:"submission_ids,omitempty"`
	}

	RejudgeSubmissionOutput struct {
		RejudgeSuccessSubmissionIds []string `json:"rejudge_success_submission_ids,omitempty"`
	}
)
