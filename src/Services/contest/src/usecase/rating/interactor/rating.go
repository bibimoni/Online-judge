package ratinginteractor

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	ratingrepo "contest/src/domain/repository/rating"
	scoreboardrepo "contest/src/domain/repository/scoreboard"
	ratingservice "contest/src/service/rating"
	ratingserviceimpl "contest/src/service/rating/impl"
	ratingusecase "contest/src/usecase/rating"
	"context"
)

type RatingInteractorImpl struct {
	ratingRepo     ratingrepo.RatingRepository
	contestRepo    contestrepo.ContestRepository
	scoreboardRepo scoreboardrepo.ScoreboardRepository
	ratingService  ratingservice.RatingService
}

func NewRatingInteractor(
	ratingRepo ratingrepo.RatingRepository,
	contestRepo contestrepo.ContestRepository,
	scoreboardRepo scoreboardrepo.ScoreboardRepository,
	ratingService ratingservice.RatingService,
) ratingusecase.RatingInteractor {
	return &RatingInteractorImpl{
		ratingRepo:     ratingRepo,
		contestRepo:    contestRepo,
		scoreboardRepo: scoreboardRepo,
		ratingService:  ratingService,
	}
}

func (r *RatingInteractorImpl) ComputeContestRatings(
	ctx context.Context,
	input *ratingusecase.ComputeContestRatingsInput,
) (*ratingusecase.ComputeContestRatingsOutput, error) {
	contest, err := r.contestRepo.GetById(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}

	if !contest.IsContestManager(input.Username) && !contest.IsAdmin(input.Username, input.UserRole) {
		return nil, common.NewForbiddenError("only contest managers or admins can compute ratings")
	}

	if !contest.HasEnded() {
		return nil, common.NewBadRequestError("contest has not ended yet; cannot compute ratings")
	}

	existing, _ := r.ratingRepo.GetRatingResultsByContest(ctx, input.ContestId)
	if len(existing) > 0 {
		return nil, common.NewConflictError("ratings have already been computed for this contest")
	}

	snapshot, err := r.scoreboardRepo.GetSnapshot(ctx, input.ContestId, domain.FinalSnapshot)
	if err != nil {
		snapshot, err = r.scoreboardRepo.GetLatestSnapshot(ctx, input.ContestId)
		if err != nil {
			return nil, common.NewBadRequestError("no scoreboard snapshot available; finalize the scoreboard first")
		}
	}

	standings := ratingserviceimpl.StandingsFromScoreboard(snapshot, contest)
	if len(standings) == 0 {
		return nil, common.NewBadRequestError("no rated participants found in the scoreboard")
	}

	usernames := make([]string, len(standings))
	for i, s := range standings {
		usernames[i] = s.Username
	}
	oldStates, err := r.ratingRepo.BatchGetUserRatingStates(ctx, usernames)
	if err != nil {
		return nil, err
	}

	results, newStates, err := r.ratingService.ComputeRatings(ctx, contest, standings, oldStates)
	if err != nil {
		return nil, err
	}

	if err := r.ratingRepo.SaveRatingResults(ctx, results); err != nil {
		return nil, err
	}

	stateList := make([]domain.UserRatingState, 0, len(newStates))
	for _, st := range newStates {
		stateList = append(stateList, st)
	}
	if err := r.ratingRepo.SaveUserRatingStates(ctx, stateList); err != nil {
		return nil, err
	}

	return &ratingusecase.ComputeContestRatingsOutput{Results: results}, nil
}

func (r *RatingInteractorImpl) GetContestRatingResults(
	ctx context.Context,
	input *ratingusecase.GetContestRatingResultsInput,
) (*ratingusecase.GetContestRatingResultsOutput, error) {
	results, err := r.ratingRepo.GetRatingResultsByContest(ctx, input.ContestId)
	if err != nil {
		return nil, err
	}
	return &ratingusecase.GetContestRatingResultsOutput{Results: results}, nil
}

func (r *RatingInteractorImpl) GetUserRatingHistory(
	ctx context.Context,
	input *ratingusecase.GetUserRatingHistoryInput,
) (*ratingusecase.GetUserRatingHistoryOutput, error) {
	history, err := r.ratingRepo.GetRatingHistory(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	return &ratingusecase.GetUserRatingHistoryOutput{History: history}, nil
}

func (r *RatingInteractorImpl) GetUserRating(
	ctx context.Context,
	input *ratingusecase.GetUserRatingInput,
) (*ratingusecase.GetUserRatingOutput, error) {
	state, err := r.ratingRepo.GetUserRatingState(ctx, input.Username)
	if err != nil {
		return nil, common.NewNotFoundError(err)
	}
	return &ratingusecase.GetUserRatingOutput{State: state}, nil
}
