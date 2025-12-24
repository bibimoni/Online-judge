package impl

import (
	"math"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
)

func CalculateTotalScore(
	problemInfo *problem.ProblemServiceGetOutput,
	testScores map[int]float64,
	testMaxScores map[int]float64,
	verdicts map[int]domain.Verdict,
	i *domain.Isolate,
) (float64, float64) {
	if len(problemInfo.TestGroups) == 0 {
		// No groups, sum all test scores
		totalScore := 0.0
		maxScore := 0.0
		for idx := range testScores {
			totalScore += testScores[idx]
			maxScore += testMaxScores[idx]
		}
		return totalScore, maxScore
	}

	groupScores := make(map[string]float64)
	totalScore := 0.0
	maxScore := 0.0

	for _, group := range problemInfo.TestGroups {
		// Check dependencies
		group := domain.TestGroup{
			Name:         group.Name,
			Scoring:      domain.ScoringType(group.Scoring),
			MaxScore:     group.MaxScore,
			TestIndices:  group.TestIndices,
			Dependencies: group.Dependencies,
		}

		if !CheckGroupDependencies(&group, groupScores) {
			i.Logger.Info().Msgf("Group %s dependencies not met, score 0", group.Name)
			groupScores[group.Name] = 0.0
			continue
		}

		groupScore := CalculateGroupScore(&group, testScores, verdicts)
		i.Logger.Info().Msgf("Group %s scored %f out of %f", group.Name, groupScore, group.MaxScore)
		groupScores[group.Name] = groupScore
		totalScore += groupScore
		maxScore += group.MaxScore
	}

	return totalScore, maxScore
}

func CalculateGroupScore(
	group *domain.TestGroup,
	testScores map[int]float64,
	verdicts map[int]domain.Verdict,
) float64 {
	switch group.Scoring {
	case domain.ScoringGroup:
		for _, idx := range group.TestIndices {
			verdict := verdicts[idx]
			if verdict != domain.ACCEPTED && verdict != domain.PARTIAL_RESULT && verdict != domain.POINTS {
				return 0.0
			}
		}
		sum := 0.0
		for _, idx := range group.TestIndices {
			sum += testScores[idx]
		}
		return sum
	case domain.ScoringSum:
		sum := 0.0
		for _, idx := range group.TestIndices {
			sum += testScores[idx]
		}
		return sum
	case domain.ScoringMin:
		minScore := math.MaxFloat64
		for _, idx := range group.TestIndices {
			if testScores[idx] < minScore {
				minScore = testScores[idx]
			}
		}
		if minScore == math.MaxFloat64 {
			return 0.0
		}
		return minScore
	}
	return 0.0
}

func CheckGroupDependencies(group *domain.TestGroup, groupScores map[string]float64) bool {
	for _, depName := range group.Dependencies {
		if groupScores[depName] == 0.0 {
			return false
		}
	}
	return true
}
