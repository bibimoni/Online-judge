package constestantinteractor

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	contestservice "contest/src/service/contest"
	contestsubmissionservice "contest/src/service/contest-submission"
	contestantusecase "contest/src/usecase/contestant"
	"context"
)

type ContestantInteractor struct {
	contestService           contestservice.ContestService
	contestrepo              contestrepo.ContestRepository
	contestsubmissionservice contestsubmissionservice.ContestSubmissionService
}

func NewContestantInteractor(
	contestService contestservice.ContestService,
	contestrepo contestrepo.ContestRepository,
	contestsubcontestsubmissionservice contestsubmissionservice.ContestSubmissionService,
) *ContestantInteractor {
	return &ContestantInteractor{
		contestService:           contestService,
		contestrepo:              contestrepo,
		contestsubmissionservice: contestsubcontestsubmissionservice,
	}
}

func (ci *ContestantInteractor) Register(ctx context.Context, input *contestantusecase.RegisterInput) (*contestantusecase.RegisterOutput, error) {
	contest, err := ci.contestrepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}

	if !contest.CanRegister(input.Username, input.RegisterType, input.Role) {
		return nil, common.NewForbiddenError("you are not allowed to register for this contest")
	}

	err = ci.contestService.AddPeople(ctx, input.ContestId, contestrepo.Contestant, input.Username, input.RegisterType)
	if err != nil {
		return nil, err
	}

	return &contestantusecase.RegisterOutput{
		Registered: true,
	}, nil
}

func (ci *ContestantInteractor) Unregister(ctx context.Context, input *contestantusecase.UnregisterInput) (*contestantusecase.UnregisterOutput, error) {
	contest, err := ci.contestrepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}

	if !contest.CanUnregister(input.Username, input.Role) {
		return nil, common.NewForbiddenError("you are not allowed to unregister from this contest")
	}

	err = ci.contestService.RemovePeople(ctx, input.ContestId, contestrepo.Contestant, input.Username)
	if err != nil {
		return nil, err
	}

	return &contestantusecase.UnregisterOutput{
		Registered: false,
	}, nil
}

func (ci *ContestantInteractor) Submit(ctx context.Context, input *contestantusecase.SubmitInput) (*contestantusecase.SubmitOutput, error) {
	contest, err := ci.contestrepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}

	if !contest.ContestantExist(input.Username) {
		return nil, common.NewForbiddenError("you are not registered in this contest")
	}

	if !contest.CanSubmit(input.Username, input.Role) {
		return nil, common.NewForbiddenError("you are not allowed to submit in this contest")
	}

	submissionId, err := ci.contestsubmissionservice.SubmitContestSubmissionToJudge(
		ctx,
		input.ContestId,
		input.ProblemLabel,
		input.Code,
		input.Language,
		input.Username,
		domain.ParticipantType(input.SubmissionType),
	)

	return &contestantusecase.SubmitOutput{
		Id:      submissionId,
		Message: "Submit successfully",
	}, err
}
