package contestantusecase

import (
	domain "contest/src/domain/entity"
	"context"
)

type ContestantInteractor interface {
	Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error)
	Unregister(ctx context.Context, input *UnregisterInput) (*UnregisterOutput, error)
	Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error)
}

type (
	RegisterInput struct {
		ContestId    string                 `json:"contest_id,omitempty"`
		RegisterType domain.ParticipantType `json:"register_type,omitempty"`
		Username     string
		Role         string
	}

	RegisterOutput struct {
		Registered bool `json:"registered"`
	}

	UnregisterInput struct {
		ContestId string `json:"contest_id,omitempty"`
		Username  string
		Role      string
	}

	UnregisterOutput struct {
		Registered bool `json:"registered"`
	}

	SubmitInput struct {
		Username       string
		Role           string
		ContestId      string         `json:"contest_id,omitempty"`
		ProblemLabel   string         `json:"problem_label,omitempty"`
		Code           string         `json:"code,omitempty"`
		Language       string         `json:"language,omitempty"`
		SubmissionType SubmissionType `json:"submission_type,omitempty"`
	}

	SubmitOutput struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}
)
type SubmissionType string

const (
	CUSTOM SubmissionType = "CUSTOM"
	ICPC   SubmissionType = "ICPC"
	IOI    SubmissionType = "IOI"
)
