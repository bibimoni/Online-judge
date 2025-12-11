package editcontest

import (
	"context"
)

type EditContestInput struct {
	EditType   string `json:"-"` // From Query Param
	ContestId  string `json:"contest-id" validate:"required"`
	PeopleType string `json:"people-type" validate:"required"`
	UserId     uint64 `json:"user-id" validate:"min=1"`
}

type EditContestOutput struct {
	Status string `json:"status"`
}

type EditContestInteractor interface {
	EditContest(ctx context.Context, input *EditContestInput) (*EditContestOutput, error)
}
