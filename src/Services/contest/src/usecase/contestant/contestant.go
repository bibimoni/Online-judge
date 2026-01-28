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
		ContestId    string
		Username     string
		RegisterType domain.ParticipantType
	}

	RegisterOutput struct {
		Registered bool `json:"registered"`
	}

	UnregisterInput struct {
		ContestId string
		Username  string
	}

	UnregisterOutput struct {
		Registered bool `json:"registered"`
	}

	SubmitInput struct {
		Username       string         `json:"username,omitempty"`
		ContestId      string         `json:"contest_id,omitempty"`
		ProblemLabel   string         `json:"problem_id,omitempty"`
		Code           string         `json:"code,omitempty"`
		LanguageId     string         `json:"language,omitempty"`
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
