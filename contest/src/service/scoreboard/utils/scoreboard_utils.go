package scoreboardserviceutils

import (
	domain "contest/src/domain/entity"
	"time"
)

// BuildScoreboardProblemResult builds the ScoreboardProblemResult for a specific problem
func BuildScoreboardProblemResult(
	ContestStartTime time.Time,
	problem domain.ContestProblem,
	isFirstToSolve bool,
	userContestProblemSubmissions []domain.ContestSubmission,
) domain.ScoreboardProblemResult {
	result := domain.ScoreboardProblemResult{
		ProblemId: problem.ProblemId,
		Label:     problem.Label,
		Points: 0.0,
		Attempts: 0,
		Solved: false,
		SolvedAt : -1,
		Pending: false,
		FirstSolve: isFirstToSolve,
		PenaltyAttempts: 0,
	}

	// Analyze submissions to build the problem result
	for i := range userContestProblemSubmissions {
		submission := userContestProblemSubmissions[i]
		result.Attempts++	
		if submission.EvalStatus == domain.Pending {
			result.Pending = true
		}
		if submission.Verdict == domain.Accepted {
			result.Solved = true
		}
		if submission.Points > result.Points {
			result.Points = submission.Points
			result.BestContestSubmissionId = submission.Id
		}

		if !result.Solved && submission.Verdict != domain.Accepted {
			result.PenaltyAttempts++
		}

		if result.Solved && result.SolvedAt == -1 {
			result.SolvedAt = int(submission.SubmitAt.Sub(ContestStartTime).Minutes())
		}
	}
	return result
}

func GetUserProblemSubmissions(
	contestSubmissions []domain.ContestSubmission,
) map[string]map[uint64][]domain.ContestSubmission {
	userProblemSubmissions := make(map[string]map[uint64][]domain.ContestSubmission)
	for i := range contestSubmissions {
		submission := contestSubmissions[i]
		if _, exists := userProblemSubmissions[submission.Username]; !exists {
			userProblemSubmissions[submission.Username] = make(map[uint64][]domain.ContestSubmission)
		}
		userProblemSubmissions[submission.Username][submission.ProblemId] = append(
			userProblemSubmissions[submission.Username][submission.ProblemId],
			submission,
		)
	}
	return userProblemSubmissions
}

func BuildFirstSolveResults(
	contestSubmissions []domain.ContestSubmission,
) map[uint64]string {
	firstSolveResults := make(map[uint64]string)
	bestProblemSolveTime := make(map[uint64]time.Time)
	for i := range contestSubmissions {
		submission := contestSubmissions[i]
		if submission.Verdict == domain.Accepted {
			if _, exists := firstSolveResults[submission.ProblemId]; !exists {
				firstSolveResults[submission.ProblemId] = submission.Username
				bestProblemSolveTime[submission.ProblemId] = submission.SubmitAt
			} else if submission.SubmitAt.Before(bestProblemSolveTime[submission.ProblemId]) {
				firstSolveResults[submission.ProblemId] = submission.Username
				bestProblemSolveTime[submission.ProblemId] = submission.SubmitAt
			}
		}
	}
	return firstSolveResults
}

func CalculateScore(
	problemResult []domain.ScoreboardProblemResult,
)(float64) {
	result := 0.0
	for i := range problemResult {
		result += problemResult[i].Points
	}
	return result
}