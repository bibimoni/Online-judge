package repository

import (
	"context"

	_ "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubmissionRepository interface {
	GetCollectionName() string
	CreateSubmission(ctx context.Context, params CreateSubmissionInput) (string, error)
	FindSubmission(ctx context.Context, submissionId string) (*domain.Submission, error)
	FindAllProblemSubmissionIds(ctx context.Context, problemId string) ([]string, error)
	GetCollection() *mongo.Collection
}

type CreateSubmissionInput struct {
	ProblemId string
	Username  string
	Type      domain.SubmissionType
}
