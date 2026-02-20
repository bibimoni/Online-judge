package scoreboardserviceimpl

import (
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	contestsubmissionrepo "contest/src/domain/repository/contest-submission"
	scoreboardrepo "contest/src/domain/repository/scoreboard"
	scoreboardservice "contest/src/service/scoreboard"
	scoreboardserviceutils "contest/src/service/scoreboard/utils"
	"context"
	"slices"
)

type ScoreboardServiceImpl struct {
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository
	contestrepo           contestrepo.ContestRepository
	scoreboardrepo        scoreboardrepo.ScoreboardRepository
}

func NewScoreboardServiceImpl(
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository,
	contestrepo contestrepo.ContestRepository,
	scoreboardrepo scoreboardrepo.ScoreboardRepository,
) *ScoreboardServiceImpl {
	return &ScoreboardServiceImpl{
		contestsubmissionrepo: contestsubmissionrepo,
		contestrepo:           contestrepo,
		scoreboardrepo:        scoreboardrepo,
	}
}

func NewScoreboardService(
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository,
	contestrepo contestrepo.ContestRepository,
	scoreboardrepo scoreboardrepo.ScoreboardRepository,
) scoreboardservice.ScoreboardService {
	return NewScoreboardServiceImpl(contestsubmissionrepo, contestrepo, scoreboardrepo)
}

func (s *ScoreboardServiceImpl) BuildScoreboardSnapshot(
	ctx context.Context,
	contestId string,
	includeVirtual bool,
	includeUnrated bool,
) (*domain.ScoreboardSnapshot, error) {
	contest, err := s.contestrepo.GetById(ctx, contestId)
	if err != nil {
		return nil, err
	}

	switch contest.ContestRule.ScoringType {
	case domain.ICPC:
		return s.ICPCSnapshot(ctx, contest, includeVirtual, includeUnrated)
	case domain.IOI:
		return s.IOISnapshot(ctx, contest, includeVirtual, includeUnrated)
	default:
		return nil, scoreboardservice.ErrInvalidScoringType
	}
}

func (c *ScoreboardServiceImpl) ICPCSnapshot(
	ctx context.Context,
	contest *domain.Contest,
	includeVirtual bool,
	includeUnrated bool,
) (*domain.ScoreboardSnapshot, error) {
	contestSubmissions, err := c.contestsubmissionrepo.ListByContest(ctx, contest.Id.Hex(), includeVirtual, includeUnrated)
	if err != nil {
		return nil, err
	}

	var snapshotKind domain.SnapshotKind
	if contest.HasEnded() {
		snapshotKind = domain.FinalSnapshot
	} else {
		snapshotKind = domain.LiveSnapshot
	}

	var userSubmission = make(map[string][]domain.ContestSubmission)
	for i := range contestSubmissions {
		submission := contestSubmissions[i]
		userSubmission[submission.Username] = append(userSubmission[submission.Username], submission)
	}

	scoreboardRows := make([]domain.ScoreboardRow, 0, len(userSubmission))
	firstSolveResults := scoreboardserviceutils.BuildFirstSolveResults(contestSubmissions)
	userProblemSubmissions := scoreboardserviceutils.GetUserProblemSubmissions(contestSubmissions)
	for username, problemSubmissions := range userProblemSubmissions {
		penalties := 0
		problemResults := make([]domain.ScoreboardProblemResult, len(contest.Problems))
		for i := range contest.Problems {
			problem := contest.Problems[i]
			isFirstSolve := firstSolveResults[problem.ProblemId] == username
			problemResults[i] = scoreboardserviceutils.BuildScoreboardProblemResult(
				contest.StartTime,
				problem,
				isFirstSolve,
				problemSubmissions[problem.ProblemId],
			)
			penalties += problemResults[i].PenaltyAttempts
		}
		scoreboardRows = append(
			scoreboardRows, domain.ScoreboardRow{
				Username: username,
				Problems: problemResults,
				Score:    scoreboardserviceutils.CalculateScore(problemResults),
				Penalty:  penalties,
			})
	}

	slices.SortFunc(scoreboardRows, func(a, b domain.ScoreboardRow) int {
		if a.Score != b.Score {
			return int(a.Score - b.Score)
		}
		return int(b.Penalty - a.Penalty)
	})

	for i := range scoreboardRows {
		scoreboardRows[i].Rank = i + 1
	}

	snapshot, err := c.scoreboardrepo.CreateAndGetSnapshot(
		ctx,
		contest.Id.Hex(),
		snapshotKind,
		scoreboardRows,
	)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
func (c *ScoreboardServiceImpl) IOISnapshot(
	ctx context.Context,
	contest *domain.Contest,
	includeVirtual bool,
	includeUnrated bool,
) (*domain.ScoreboardSnapshot, error) {
	contestSubmissions, err := c.contestsubmissionrepo.ListByContest(ctx, contest.Id.Hex(), includeVirtual, includeUnrated)
	if err != nil {
		return nil, err
	}
	var snapshotKind domain.SnapshotKind
	if contest.HasEnded() {
		snapshotKind = domain.FinalSnapshot
	} else {
		snapshotKind = domain.LiveSnapshot
	}

	var userSubmission = make(map[string][]domain.ContestSubmission)
	for i := range contestSubmissions {
		submission := contestSubmissions[i]
		userSubmission[submission.Username] = append(userSubmission[submission.Username], submission)
	}

	userProblemSubmissions := scoreboardserviceutils.GetUserProblemSubmissions(contestSubmissions)
	scoreboardRows := make([]domain.ScoreboardRow, 0, len(userSubmission))

	for username := range userSubmission {
		var problemResults []domain.ScoreboardProblemResult
		for _, problem := range contest.Problems {
			userProblemSubs := userProblemSubmissions[username][problem.ProblemId]
			problemResult := scoreboardserviceutils.BuildScoreboardProblemResult(
				contest.StartTime,
				problem,
				false,
				userProblemSubs,
			)
			problemResults = append(problemResults, problemResult)
		}
		scoreboardRows = append(
			scoreboardRows, domain.ScoreboardRow{
				Username: username,
				Score:    scoreboardserviceutils.CalculateScore(problemResults),
				Problems: problemResults,
			})
	}
	slices.SortFunc(scoreboardRows, func(a, b domain.ScoreboardRow) int {
		return int(a.Score - b.Score)
	})

	for i := range scoreboardRows {
		scoreboardRows[i].Rank = i + 1
	}

	snapshot, err := c.scoreboardrepo.CreateAndGetSnapshot(
		ctx,
		contest.Id.Hex(),
		snapshotKind,
		scoreboardRows,
	)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}
