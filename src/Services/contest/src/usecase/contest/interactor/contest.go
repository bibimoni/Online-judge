package contestinteractor

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	"contest/src/domain/repository/contest"
	"contest/src/infrastructure/config"
	contestservice "contest/src/service/contest"
	contestusecase "contest/src/usecase/contest"
	"context"
	"fmt"
	"slices"
)

type ContestInteractor struct {
	contestRepo    contestrepo.ContestRepository
	contestService contestservice.ContestService
}

func NewContestInteractor(
	contestRepo contestrepo.ContestRepository,
	contestService contestservice.ContestService,
) *ContestInteractor {
	return &ContestInteractor{
		contestRepo:    contestRepo,
		contestService: contestService,
	}
}

func (i *ContestInteractor) CreateContest(ctx context.Context, input *contestusecase.CreateContestInput) (*contestusecase.CreateContestOutput, error) {
	if !i.contestRepo.CanCreateContest(ctx, input.UserRole) {
		return nil, common.NewForbiddenError("you don't have permission to create contest")
	}
	contestId, err := i.contestService.Create(ctx, input.Username, input.Name)
	if err != nil {
		return nil, err
	}

	return &contestusecase.CreateContestOutput{
		ContestId: contestId,
	}, nil
}

func (i *ContestInteractor) EditContest(ctx context.Context, input *contestusecase.EditContestInput) (*contestusecase.EditContestOutput, error) {
	// This queries the contest twice, may need optimization later
	contest, err := i.contestRepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, common.NewNotFoundError(err)
	}

	if !contest.IsAdmin(input.Username, input.UserRole) {
		return nil, common.NewForbiddenError("only admins can edit")
	}

	if input.EditType == contestusecase.AddPeople {
		return i.handleAddPeople(ctx, input)
	}
	if input.EditType == contestusecase.RemovePeople {
		return i.handleRemovePeople(ctx, input)
	}

	return nil, fmt.Errorf("invalid edit_type: %s", input.EditType)
}

func (i *ContestInteractor) handleAddPeople(ctx context.Context, input *contestusecase.EditContestInput) (*contestusecase.EditContestOutput, error) {
	err := i.contestRepo.AddPeople(ctx, input.ContestId, input.PeopleType, input.Target)
	if err != nil {
		return nil, err
	}
	return &contestusecase.EditContestOutput{StatusOK: common.StatusOK{Status: "ok"}}, nil
}

func (i *ContestInteractor) handleRemovePeople(ctx context.Context, input *contestusecase.EditContestInput) (*contestusecase.EditContestOutput, error) {
	err := i.contestService.RemovePeople(ctx, input.ContestId, input.PeopleType, input.Target)
	if err != nil {
		return nil, err
	}
	return &contestusecase.EditContestOutput{StatusOK: common.StatusOK{Status: "ok"}}, nil
}

func (i *ContestInteractor) PatchContest(
	ctx context.Context,
	input *contestusecase.PatchContestInput,
) (*contestusecase.PatchContestOutput, error) {
	contest, err := i.contestRepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, common.NewNotFoundError(err)
	}

	if !contest.IsAdmin(input.Username, input.UserRole) {
		return nil, common.NewForbiddenError("only contest manager can make changes to this contest")
	}

	updateData := make(map[string]any)
	var (
		startTime        = contest.StartTime
		endTime          = contest.EndTime
		rejudgeWindowEnd = contest.RejudgeWindowEnd
		finalizeAt       = contest.FinalizeAt
	)

	if input.Description != nil {
		updateData["description"] = *input.Description
	}
	if input.ScoreboardVisibility != nil {
		allowedVisibilities := []domain.ScoreboardVisibility{
			domain.ScoreboardHidden,
			domain.ScoreboardPublic,
			domain.ScoreboardContestantOnly,
		}
		if !slices.Contains(allowedVisibilities, *input.ScoreboardVisibility) {
			return nil, common.NewBadRequestError("invalid scoreboard visibility")
		}
		updateData["scoreboard_visibility"] = *input.ScoreboardVisibility
	}
	if input.StartTime != nil {
		updateData["start_time"] = *input.StartTime
		startTime = *input.StartTime
	}
	if input.EndTime != nil {
		updateData["end_time"] = *input.EndTime
		endTime = *input.EndTime
	}
	if input.RejudgeWindowEnd != nil {
		updateData["rejudge_window_end"] = *input.RejudgeWindowEnd
		rejudgeWindowEnd = *input.RejudgeWindowEnd
	}
	if input.FinalizeAt != nil {
		updateData["finalize_at"] = *input.FinalizeAt
		finalizeAt = *input.FinalizeAt
	}

	// contest time validation
	if startTime.After(endTime) {
		return nil, common.NewBadRequestError("start time must sooner than end time")
	}

	// rejudge window end must sooner than finalize at
	if rejudgeWindowEnd.After(finalizeAt) {
		return nil, common.NewBadRequestError("rejudge window end must sooner than finalize at")
	}

	if input.ContestRule != nil {
		if input.ContestRule.ScoringType != nil {
			allowedScoringTypes := []domain.ScoringType{
				domain.ICPC,
				domain.IOI,
			}
			if !slices.Contains(allowedScoringTypes, *input.ContestRule.ScoringType) {
				return nil, common.NewBadRequestError("invalid scoring type")
			}
			updateData["contest_rule.scoring_type"] = *input.ContestRule.ScoringType
		}
		if input.ContestRule.PenaltyMinutes != nil {
			updateData["contest_rule.penalty_minutes"] = *input.ContestRule.PenaltyMinutes
		}
		if input.ContestRule.FreezeStartTime != nil {
			updateData["contest_rule.freeze_start_time"] = *input.ContestRule.FreezeStartTime
		}
		if input.ContestRule.FreezeTime != nil {
			updateData["contest_rule.freeze_time"] = *input.ContestRule.FreezeTime
		}
		if input.ContestRule.MaxAllowedSubmissionsPerProblem != nil {
			updateData["contest_rule.max_allowed_submissions_per_problem"] = *input.ContestRule.MaxAllowedSubmissionsPerProblem
		}
	}

	if len(updateData) == 0 {
		return &contestusecase.PatchContestOutput{StatusOK: common.StatusOK{Status: "ok"}}, nil
	}

	err = i.contestRepo.UpdateOne(ctx, input.ContestId, updateData)
	if err != nil {
		return nil, err
	}

	return &contestusecase.PatchContestOutput{StatusOK: common.StatusOK{Status: "ok"}}, nil
}

func (i *ContestInteractor) ManageContestProblems(
	ctx context.Context,
	input *contestusecase.ManageContestProblemsInput,
) (*contestusecase.ManageContestProblemsOutput, error) {
	err := i.contestService.ChangeProblems(ctx, input.ContestId, input.ProblemIds, input.ShortNames)
	if err != nil {
		return nil, err
	}
	return &contestusecase.ManageContestProblemsOutput{StatusOK: common.StatusOK{Status: "ok"}}, nil
}

func (i *ContestInteractor) GetContestById(ctx context.Context, input *contestusecase.GetContestByIdInput) (*contestusecase.GetContestByIdOutput, error) {
	contest, err := i.contestRepo.GetById(ctx, input.ContestId)
	if contest.Visibility == domain.VisibilityPublic {
		return &contestusecase.GetContestByIdOutput{
			Contest: contest,
		}, err
	}
	if !input.Authenticated {
		return nil, common.NewForbiddenError("you must be logged in to access this contest")
	}

	if contest.CanViewContest(input.Username, input.Role) {
		return &contestusecase.GetContestByIdOutput{
			Contest: contest,
		}, err
	}

	return nil, common.NewForbiddenError("you don't have permission to access this contest")
}
func (i *ContestInteractor) GetAllContests(ctx context.Context, input *contestusecase.GetContestsInput) (*contestusecase.GetContestsOutput, error) {
	config.GetLogger().Info().Msgf("GetAllContests called with input: %+v", input)
	if !input.Authenticated {
		contests, err := i.contestRepo.ListAllPublicContests(ctx)
		if err != nil {
			return nil, err
		}
		return &contestusecase.GetContestsOutput{
			Contests: contests,
		}, nil
	}
	contests, err := i.contestRepo.ListAllContestsWithAuth(ctx, input.Username, input.Role)
	if err != nil {
		return nil, err
	}
	return &contestusecase.GetContestsOutput{
		Contests: contests,
	}, nil
}
