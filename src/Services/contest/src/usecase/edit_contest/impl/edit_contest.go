package editcontestimpl

import (
	repository "contest/src/domain/repository/contest"
	editcontest "contest/src/usecase/edit_contest"
	"context"
	"fmt"
)

type editContestInteractor struct {
	contestRepo repository.ContestRepository
}

func NewEditContestInteractor(contestRepo repository.ContestRepository) editcontest.EditContestInteractor {
	return &editContestInteractor{
		contestRepo: contestRepo,
	}
}

func (i *editContestInteractor) EditContest(ctx context.Context, input *editcontest.EditContestInput) (*editcontest.EditContestOutput, error) {
	if input.EditType == "add-people" {
		return i.handleAddPeople(ctx, input)
	}
	if input.EditType == "remove-people" {
		return i.handleRemovePeople(ctx, input)
	}

	return nil, fmt.Errorf("invalid edit-type: %s", input.EditType)
}

func (i *editContestInteractor) handleAddPeople(ctx context.Context, input *editcontest.EditContestInput) (*editcontest.EditContestOutput, error) {
	err := i.contestRepo.AddPeople(input.ContestId, input.PeopleType, input.UserId)
	if err != nil {
		return nil, err
	}
	return &editcontest.EditContestOutput{Status: "ok"}, nil
}

func (i *editContestInteractor) handleRemovePeople(ctx context.Context, input *editcontest.EditContestInput) (*editcontest.EditContestOutput, error) {
	err := i.contestRepo.RemovePeople(input.ContestId, input.PeopleType, input.UserId)
	if err != nil {
		return nil, err
	}
	return &editcontest.EditContestOutput{Status: "ok"}, nil
}
