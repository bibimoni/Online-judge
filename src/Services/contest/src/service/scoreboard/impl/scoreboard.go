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
	"time"

	"github.com/rs/zerolog/log"
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

	// Optimization: if the contest has ended, try to return the existing final snapshot
	// from the database instead of rebuilding every time.
	// Only rebuild if there are new submissions since the last snapshot
	// (e.g. from virtual or unrated participants).
	if contest.HasEnded() {
		existingSnapshot, err := s.scoreboardrepo.GetSnapshot(ctx, contestId, domain.FinalSnapshot)
		if err == nil && existingSnapshot != nil {
			// Fetch all submissions and check if any arrived after the snapshot was created
			allSubmissions, listErr := s.contestsubmissionrepo.ListByContest(ctx, contestId, includeVirtual, includeUnrated)
			if listErr != nil {
				log.Warn().Err(listErr).Msg("Failed to list submissions, will rebuild scoreboard")
			} else {
				hasNewSubmissions := false
				for i := range allSubmissions {
					if allSubmissions[i].SubmitAt.After(existingSnapshot.CreatedAt) ||
						allSubmissions[i].UpdatedAt.After(existingSnapshot.CreatedAt) {
						hasNewSubmissions = true
						break
					}
				}
				if !hasNewSubmissions {
					// No new submissions — return cached snapshot directly
					log.Debug().Msgf("Returning cached final scoreboard for contest %s (no new submissions)", contestId)
					return existingSnapshot, nil
				}
				log.Info().Msgf("Found new submissions after last snapshot for contest %s, rebuilding", contestId)
			}
		}
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

// GetScoreboardAt returns the scoreboard as it would be at the given point in time.
// It filters submissions to only include those submitted at or before the given time,
// computes the scoreboard in-memory, and returns it without persisting to the database.
func (s *ScoreboardServiceImpl) GetScoreboardAt(
	ctx context.Context,
	contestId string,
	at time.Time,
	includeVirtual bool,
	includeUnrated bool,
) (*domain.ScoreboardSnapshot, error) {
	contest, err := s.contestrepo.GetById(ctx, contestId)
	if err != nil {
		return nil, err
	}

	// Fetch all submissions, then filter to only those submitted at or before the given time
	allSubmissions, err := s.contestsubmissionrepo.ListByContest(ctx, contestId, includeVirtual, includeUnrated)
	if err != nil {
		return nil, err
	}

	contestSubmissions := filterSubmissionsBefore(allSubmissions, at)

	switch contest.ContestRule.ScoringType {
	case domain.ICPC:
		return s.buildICPCScoreboardFromSubmissions(contest, contestSubmissions), nil
	case domain.IOI:
		return s.buildIOIScoreboardFromSubmissions(contest, contestSubmissions), nil
	default:
		return nil, scoreboardservice.ErrInvalidScoringType
	}
}

// buildICPCScoreboardFromSubmissions builds an ICPC scoreboard from a given set of submissions
// without persisting to the database.
func (s *ScoreboardServiceImpl) buildICPCScoreboardFromSubmissions(
	contest *domain.Contest,
	contestSubmissions []domain.ContestSubmission,
) *domain.ScoreboardSnapshot {
	firstSolveResults := scoreboardserviceutils.BuildFirstSolveResults(contestSubmissions)
	userProblemSubmissions := scoreboardserviceutils.GetUserProblemSubmissions(contestSubmissions)

	scoreboardRows := make([]domain.ScoreboardRow, 0, len(userProblemSubmissions))
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

	var snapshotKind domain.SnapshotKind
	if contest.HasEnded() {
		snapshotKind = domain.FinalSnapshot
	} else {
		snapshotKind = domain.LiveSnapshot
	}

	return &domain.ScoreboardSnapshot{
		ContestId: contest.Id,
		Kind:      snapshotKind,
		Rows:      scoreboardRows,
		CreatedAt: time.Now(),
	}
}

// buildIOIScoreboardFromSubmissions builds an IOI scoreboard from a given set of submissions
// without persisting to the database.
func (s *ScoreboardServiceImpl) buildIOIScoreboardFromSubmissions(
	contest *domain.Contest,
	contestSubmissions []domain.ContestSubmission,
) *domain.ScoreboardSnapshot {
	userProblemSubmissions := scoreboardserviceutils.GetUserProblemSubmissions(contestSubmissions)

	// Collect unique usernames
	userSubmission := make(map[string]bool)
	for i := range contestSubmissions {
		userSubmission[contestSubmissions[i].Username] = true
	}

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
		scoreboardRows = append(scoreboardRows, domain.ScoreboardRow{
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

	var snapshotKind domain.SnapshotKind
	if contest.HasEnded() {
		snapshotKind = domain.FinalSnapshot
	} else {
		snapshotKind = domain.LiveSnapshot
	}

	return &domain.ScoreboardSnapshot{
		ContestId: contest.Id,
		Kind:      snapshotKind,
		Rows:      scoreboardRows,
		CreatedAt: time.Now(),
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

	snapshot := c.buildICPCScoreboardFromSubmissions(contest, contestSubmissions)

	// Persist snapshot to DB
	persistedSnapshot, err := c.scoreboardrepo.CreateAndGetSnapshot(
		ctx,
		contest.Id.Hex(),
		snapshot.Kind,
		snapshot.Rows,
	)
	if err != nil {
		return nil, err
	}

	return persistedSnapshot, nil
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

	snapshot := c.buildIOIScoreboardFromSubmissions(contest, contestSubmissions)

	// Persist snapshot to DB
	persistedSnapshot, err := c.scoreboardrepo.CreateAndGetSnapshot(
		ctx,
		contest.Id.Hex(),
		snapshot.Kind,
		snapshot.Rows,
	)
	if err != nil {
		return nil, err
	}

	return persistedSnapshot, nil
}

// filterSubmissionsBefore returns only submissions with SubmitAt at or before the given time.
func filterSubmissionsBefore(submissions []domain.ContestSubmission, before time.Time) []domain.ContestSubmission {
	filtered := make([]domain.ContestSubmission, 0, len(submissions))
	for i := range submissions {
		if !submissions[i].SubmitAt.After(before) {
			filtered = append(filtered, submissions[i])
		}
	}
	return filtered
}
