package contestusecase

import "context"

type ContestOutput struct {
	Status string `json:"status"`
}

type ContestInteractor interface {
	EditContest(ctx context.Context, input *EditContestInput) (*EditContestOutput, error)
	CreateContest(ctx context.Context, input *CreateContestInput) (*CreateContestOutput, error)
}

type (
	EditContestInput struct {
		EditType   string `json:"-"` // From Query Param
		ContestId  string `json:"contest-id" validate:"required"`
		PeopleType string `json:"people-type" validate:"required"`
		Username   string `json:"username" validate:"min=6"`
	}

	EditContestOutput struct {
		Status string `json:"status"`
	}

	CreateContestInput struct {
		Creator string
	}

	CreateContestOutput struct {
		ContestId string `json:"contest-id"`
	}
)
