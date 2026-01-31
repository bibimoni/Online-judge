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

type ICPCCalculator struct{
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository
	contestrepo contestrepo.ContestRepository
	scoreboardrepo scoreboardrepo.ScoreboardRepository
}

func NewICPCCalculator(
	contestsubmissionrepo contestsubmissionrepo.ContestSubmissionRepository,
	 contestrepo contestrepo.ContestRepository,
	  scoreboardrepo scoreboardrepo.ScoreboardRepository,
) *ICPCCalculator {
	return &ICPCCalculator{
		contestsubmissionrepo: contestsubmissionrepo,
		contestrepo: contestrepo,
		scoreboardrepo: scoreboardrepo,	
	}
}

func (c *ICPCCalculator) GetScoringType() domain.ScoringType {
	return domain.ICPC
}

func (c *ICPCCalculator) BuildScoreboardSnapshot(
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
		scoreboardRows = append(scoreboardRows, domain.ScoreboardRow{
			Username:       username,
			Problems: problemResults,
			Score: scoreboardserviceutils.CalculateScore(problemResults),
			Penalty:      penalties,
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
		contestId,
		snapshotKind,
		scoreboardRows,
	) 
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}