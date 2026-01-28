package contestsubmissionrepoimpl

import (
	domain "contest/src/domain/entity"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ContestSubmissionRepositoryImpl struct {
	collection *mongo.Collection
}

func NewContestSubmissionRepositoryImpl(db *mongo.Database) *ContestSubmissionRepositoryImpl {
	return &ContestSubmissionRepositoryImpl{
		collection: db.Collection("ContestSubmission"),
	}
}

func NewContestSubmissionRepository(db *mongo.Database) *ContestSubmissionRepositoryImpl {
	return NewContestSubmissionRepositoryImpl(db)
}

func (csr *ContestSubmissionRepositoryImpl) Create(
	ctx context.Context,
	contestId string,
	username string,
	contestProblemId string,
	submissionType domain.ParticipantType,
	submitAt time.Time,
	submissionId string,
) (string, error) {
	contestBsonId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return "", err
	}

	contestProblemBsonId, err := bson.ObjectIDFromHex(contestProblemId)
	newContestSubmission := domain.ContestSubmission{
		ContestId:        contestBsonId,
		SubmissionId:     submissionId,
		Username:         username,
		ContestProblemId: contestProblemBsonId,
		SubmitAt:         submitAt,
		Points:           0,
		EvalStatus:       domain.Pending,
		Ignored:          false,
		UpdatedAt:        submitAt,
		SubmissionType:   submissionType,
	}

	result, err := csr.collection.InsertOne(ctx, newContestSubmission)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (csr *ContestSubmissionRepositoryImpl) UpsertFromJudgeEvent(
	ctx context.Context,
	submissionId string,
	verdict domain.Verdict,
	points float64,
) (string, error) {
	result, err := csr.collection.UpdateOne(ctx,
		bson.M{"submission_id": submissionId},
		bson.M{
			"$set": bson.M{
				"submission_id": submissionId,
				"eval_status":   domain.Finished,
				"verdict":       verdict,
				"points":        points,
				"updated_at":    time.Now(),
			},
		})

	if err != nil {
		return "", err
	}

	return result.UpsertedID.(bson.ObjectID).Hex(), nil
}

func (csr *ContestSubmissionRepositoryImpl) ListByContest(
	ctx context.Context,
	contestId string,
	includeVirtual bool,
	includeUnrated bool,
) ([]domain.ContestSubmission, error) {
	contestBsonId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return nil, err
	}

	cursor, err := csr.collection.Find(ctx,
		bson.M{
			"contest_id": contestBsonId,
		},
	)

	if err != nil {
		return nil, err
	}

	var contestSubmissions []domain.ContestSubmission
	if err := cursor.All(ctx, &contestSubmissions); err != nil {
		return nil, err
	}

	return contestSubmissions, nil
}

func (csr *ContestSubmissionRepositoryImpl) SetIgnored(ctx context.Context, contestSubmissionId string, ignored bool) error {
	cId, err := bson.ObjectIDFromHex(contestSubmissionId)
	if err != nil {
		return err
	}

	_, err = csr.collection.UpdateOne(ctx, bson.M{"_id": cId}, bson.M{"$set": bson.M{"ignored": ignored}})
	if err != nil {
		return err
	}

	return nil
}

func (csr *ContestSubmissionRepositoryImpl) SetIgnoreUser(ctx context.Context, contestId string, username string, ignored bool) error {
	contestBsonId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return err
	}

	_, err = csr.collection.UpdateMany(ctx,
		bson.M{
			"contest_id": contestBsonId,
			"username":   username,
		},
		bson.M{
			"$set": bson.M{
				"ignored": ignored,
			},
		},
	)

	return err
}
