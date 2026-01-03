package contestinteractor

import (
	"contest/src/domain/repository/contest"
	contestusecase "contest/src/usecase/contest"
	"context"
	"fmt"
)

type ContestInteractor struct {
	contestRepo contestrepo.ContestRepository
}

func NewContestInteractor(contestRepo contestrepo.ContestRepository) *ContestInteractor {
	return &ContestInteractor{
		contestRepo: contestRepo,
	}
}

func (i *ContestInteractor) CreateContest(ctx context.Context, input *contestusecase.CreateContestInput) (*contestusecase.CreateContestOutput, error) {
	contestId, err := i.contestRepo.Create(ctx, input.Creator)
	if err != nil {
		return nil, err
	}

	return &contestusecase.CreateContestOutput{
		ContestId: contestId,
	}, nil
}

func (i *ContestInteractor) EditContest(ctx context.Context, input *contestusecase.EditContestInput) (*contestusecase.EditContestOutput, error) {
	if input.EditType == "add-people" {
		return i.handleAddPeople(ctx, input)
	}
	if input.EditType == "remove-people" {
		return i.handleRemovePeople(ctx, input)
	}

	return nil, fmt.Errorf("invalid edit-type: %s", input.EditType)
}

func (i *ContestInteractor) handleAddPeople(ctx context.Context, input *contestusecase.EditContestInput) (*contestusecase.EditContestOutput, error) {
	err := i.contestRepo.AddPeople(ctx, input.ContestId, input.PeopleType, input.Username)
	if err != nil {
		return nil, err
	}
	return &contestusecase.EditContestOutput{Status: "ok"}, nil
}

func (i *ContestInteractor) handleRemovePeople(ctx context.Context, input *contestusecase.EditContestInput) (*contestusecase.EditContestOutput, error) {
	err := i.contestRepo.RemovePeople(ctx, input.ContestId, input.PeopleType, input.Username)
	if err != nil {
		return nil, err
	}
	return &contestusecase.EditContestOutput{Status: "ok"}, nil
}
