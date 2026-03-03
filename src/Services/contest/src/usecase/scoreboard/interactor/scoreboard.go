package scoreboardineractor

import (
	"contest/src/common"
	contestrepo "contest/src/domain/repository/contest"
	scoreboardrepo "contest/src/domain/repository/scoreboard"
	scoreboardservice "contest/src/service/scoreboard"
	scoreboardusecase "contest/src/usecase/scoreboard"
	"context"
)

type ScoreboardInteractor struct {
	scoreboardRepo    scoreboardrepo.ScoreboardRepository
	scoreboardService scoreboardservice.ScoreboardService
	contestRepo 	 contestrepo.ContestRepository
}

func NewScoreboardInteractor(
	scoreboardRepo scoreboardrepo.ScoreboardRepository,
	scoreboardService scoreboardservice.ScoreboardService,
	contestRepo contestrepo.ContestRepository,
) *ScoreboardInteractor {
	return &ScoreboardInteractor{
		scoreboardRepo:    scoreboardRepo,
		scoreboardService: scoreboardService,
		contestRepo:       contestRepo,
	}
}

func (s *ScoreboardInteractor) BuildScoreboardSnapshot(
	ctx context.Context,
	input *scoreboardusecase.BuildScoreboardSnapshotInput,
) (*scoreboardusecase.BuildScoreboardSnapshotOutput, error) {
	contest, err := s.contestRepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}
	if !input.Authenticated && !contest.HasStarted() {
		return nil, common.NewForbiddenError("contest has not started yet")
	}
	if input.Authenticated && !contest.CanViewScoreboard(input.Username, input.UserRole) {
		return nil, common.NewForbiddenError("you do not have permission to view the scoreboard")
	}
	snapshot, err := s.scoreboardService.BuildScoreboardSnapshot(
		ctx,
		input.ContestId,
		input.IncludeVirtual,
		input.IncludeUnrated,
	)
	if err != nil {
		return nil, err
	}
	return &scoreboardusecase.BuildScoreboardSnapshotOutput{
		ScoreboardSnapshot: snapshot,
	}, nil
}

func (s *ScoreboardInteractor) GetScoreboardAt(
	ctx context.Context,
	input *scoreboardusecase.GetScoreboardAtInput,
) (*scoreboardusecase.BuildScoreboardSnapshotOutput, error) {
	contest, err := s.contestRepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}
	if !input.Authenticated && !contest.HasStarted() {
		return nil, common.NewForbiddenError("contest has not started yet")
	}
	if input.Authenticated && !contest.CanViewScoreboard(input.Username, input.UserRole) {
		return nil, common.NewForbiddenError("you do not have permission to view the scoreboard")
	}

	// Clamp the requested time: must be >= contest start time
	at := input.At
	if at.Before(contest.StartTime) {
		at = contest.StartTime
	}

	snapshot, err := s.scoreboardService.GetScoreboardAt(
		ctx,
		input.ContestId,
		at,
		input.IncludeVirtual,
		input.IncludeUnrated,
	)
	if err != nil {
		return nil, err
	}
	return &scoreboardusecase.BuildScoreboardSnapshotOutput{
		ScoreboardSnapshot: snapshot,
	}, nil
}