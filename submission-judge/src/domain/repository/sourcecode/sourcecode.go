package repository

import (
	"context"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// This will save source code into mongoDB, this is to save user's submission
type SourcecodeRepository interface {
	GetCollectionName() string
	CreateSourcecode(ctx context.Context, source string, languageId string, submissionId string) (string, error)
	GetSourcecode(ctx context.Context, id string) (*domain.SourceCode, error)
	GetSourcecodeBson(ctx context.Context, bid bson.ObjectID) (*domain.SourceCode, error)
	GetSourceBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.SourceCode, error)
}
