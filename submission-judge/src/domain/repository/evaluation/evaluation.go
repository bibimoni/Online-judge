package repository

import (
	"context"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type EvaluationRepository interface {
	GetCollectionName() string
	CreateEval(ctx context.Context, submissionId string, TL int, ML memory.Memory, nCase int) (string, error)
	UpdateVerdict(ctx context.Context, evalId string, vert domain.Verdict) error
	UpdateCase(ctx context.Context, evalId string, verdictCase domain.Verdict, cpuTimeCase float64, memoryUsageCase memory.Memory, outputCase string, pointsCase int, cpuTime float64, memoryUsage memory.Memory, nsucess int) error
	UpdateFinal(ctx context.Context, evalId string, verdict domain.Verdict, cpuTime float64, memoryUsage memory.Memory, nsucess int, points int, message string) error
	UpdateCaseFloat(ctx context.Context, evalId string, verdictCase domain.Verdict, cpuTimeCase float64, memoryUsageCase memory.Memory, outputCase string, score float64, maxScore float64, cpuTime float64, memoryUsage memory.Memory, nsucess int) error
	UpdateFinalFloat(ctx context.Context, evalId string, verdict domain.Verdict, cpuTime float64, memoryUsage memory.Memory, nsucess int, totalScore float64, totalMaxScore float64, message string) error
	GetEval(ctx context.Context, evalId string) (*domain.EvaluationResult, error)
	GetEvalBson(ctx context.Context, evalId bson.ObjectID) (*domain.EvaluationResult, error)
	GetEvalBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.EvaluationResult, error)
	GetEvalBySubmissionIdNoBson(ctx context.Context, submissionId string) (*domain.EvaluationResult, error)
	// FilterJudgedSubmissions(ctx context.Context, submissionIds []string) ([]string, error)
}
