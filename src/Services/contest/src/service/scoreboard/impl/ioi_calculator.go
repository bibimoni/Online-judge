package scoreboardserviceimpl

import (
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	contestsubmissionrepo "contest/src/domain/repository/contest-submission"
	scoreboardrepo "contest/src/domain/repository/scoreboard"
	scoreboardserviceutils "contest/src/service/scoreboard/utils"
	"context"
	"slices"
)
type IOICalculator struct{
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository
	contestrepo contestrepo.ContestRepository
	scoreboardrepo scoreboardrepo.ScoreboardRepository
}

func NewIOICalculator(
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository,
	 contestrepo contestrepo.ContestRepository,
	  scoreboardrepo scoreboardrepo.ScoreboardRepository,
) *IOICalculator {
	return &IOICalculator{
		contestsubmissionrepo: contestsubmissionrepo,
		contestrepo: contestrepo,
		scoreboardrepo: scoreboardrepo,	
	}
}

func (c *IOICalculator) GetScoringType() domain.ScoringType {
	return domain.IOI	
}

func (c *IOICalculator) BuildScoreboardSnapshot(
	ctx context.Context,
	contestId string,
	includeVirtual bool,
 	includeUnrated bool,
) (*domain.ScoreboardSnapshot, error) {
	contestSubmissions, err := c.contestsubmissionrepo.ListByContest(ctx, contestId, includeVirtual, includeUnrated)
	if err != nil {
		return nil, err
	}
	contest, err := c.contestrepo.GetById(ctx, contestId)
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
		contestId,
		snapshotKind,
		scoreboardRows,
	)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}