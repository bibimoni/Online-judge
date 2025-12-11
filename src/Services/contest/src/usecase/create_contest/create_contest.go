package createcontest

import (
	"context"
)

type CreateContestInput struct {
	AuthorId uint64
}

type CreateContestOutput struct {
	ContestId string `json:"contest-id"`
}

type CreateContestInteractor interface {
	CreateContest(ctx context.Context, input *CreateContestInput) (*CreateContestOutput, error)
}
