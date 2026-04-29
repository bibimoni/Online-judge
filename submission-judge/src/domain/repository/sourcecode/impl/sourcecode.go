package impl

import (
	"context"
	"fmt"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	repository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SourcecodeRepositoryImpl struct {
	collectionName string
	collection     *mongo.Collection
}

func NewSourcecodeRepositoryImpl(db *mongo.Database) *SourcecodeRepositoryImpl {
	collectionName := "Sourcecode"
	return &SourcecodeRepositoryImpl{
		collectionName: collectionName,
		collection:     db.Collection(collectionName),
	}
}

func NewSourcecodeRepository(db *mongo.Database) repository.SourcecodeRepository {
	return NewSourcecodeRepositoryImpl(db)
}

func (sr *SourcecodeRepositoryImpl) GetCollectionName() string {
	return sr.collectionName
}

func (sr *SourcecodeRepositoryImpl) CreateSourcecode(ctx context.Context, source string, languageId string, submissionId string) (string, error) {
	if !store.DefaultStore.Contains(languageId) {
		return "", fmt.Errorf("currently there is no support for language %s", languageId)
	}

	sid, err := bson.ObjectIDFromHex(submissionId)
	if err != nil {
		return "", err
	}
	newSource := domain.SourceCode{
		LanguageId:   languageId,
		CreatedAt:    time.Now(),
		FileSize:     memory.Memory(len(source)),
		SourceCode:   source,
		SubmissionId: sid,
	}

	got, err := sr.collection.InsertOne(ctx, newSource)

	if err != nil {
		return "", nil
	}
	return got.InsertedID.(bson.ObjectID).Hex(), nil
}

func (sr *SourcecodeRepositoryImpl) GetSourcecode(ctx context.Context, id string) (*domain.SourceCode, error) {
	var returnSourceCode domain.SourceCode
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = sr.collection.FindOne(ctx, bson.M{"_id": bid}).Decode(&returnSourceCode)
	if err != nil {
		return nil, err
	}
	return &returnSourceCode, nil
}

func (sr *SourcecodeRepositoryImpl) GetSourcecodeBson(ctx context.Context, bid bson.ObjectID) (*domain.SourceCode, error) {
	var returnSourceCode domain.SourceCode
	err := sr.collection.FindOne(ctx, bson.M{"_id": bid}).Decode(&returnSourceCode)
	if err != nil {
		return nil, err
	}
	return &returnSourceCode, nil
}

func (sr *SourcecodeRepositoryImpl) GetSourceBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.SourceCode, error) {
	var returnSourceCode domain.SourceCode
	err := sr.collection.FindOne(ctx, bson.M{"submission_id": submissionId}).Decode(&returnSourceCode)
	if err != nil {
		return nil, err
	}
	return &returnSourceCode, nil
}
