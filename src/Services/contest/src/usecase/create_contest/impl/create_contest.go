package createcontestimpl

import (
	repository "contest/src/domain/repository/contest"
	createcontest "contest/src/usecase/create_contest"
	"context"
)

type createContestInteractor struct {
	contestRepo repository.ContestRepository
}

func NewCreateContestInteractor(contestRepo repository.ContestRepository) createcontest.CreateContestInteractor {
	return &createContestInteractor{
		contestRepo: contestRepo,
	}
}

func (i *createContestInteractor) CreateContest(ctx context.Context, input *createcontest.CreateContestInput) (*createcontest.CreateContestOutput, error) {
	contestId, err := i.contestRepo.Create(ctx, input.AuthorId)
	if err != nil {
		return nil, err
	}

	return &createcontest.CreateContestOutput{
		ContestId: contestId,
	}, nil
}
