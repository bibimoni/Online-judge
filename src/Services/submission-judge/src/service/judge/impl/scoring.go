package impl

import (
	"math"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
)

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
