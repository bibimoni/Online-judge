package impl

import (
	"context"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	repository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EvaluationRepositoryImpl struct {
	collectionName string
	collection     *mongo.Collection
}

func NewEvaluationRepositoryImpl(db *mongo.Database) *EvaluationRepositoryImpl {
	collectionName := "Evaluation"
	return &EvaluationRepositoryImpl{
		collectionName: collectionName,
		collection:     db.Collection(collectionName),
	}
}

func NewEvaluationRepository(db *mongo.Database) repository.EvaluationRepository {
	return NewEvaluationRepositoryImpl(db)
}

func (er *EvaluationRepositoryImpl) GetCollectionName() string {
	return er.collectionName
}

func (er *EvaluationRepositoryImpl) CreateEval(ctx context.Context, submissionId string, TL int, ML memory.Memory, nCase int) (string, error) {
	sid, err := bson.ObjectIDFromHex(submissionId)
	if err != nil {
		return "", err
	}

	newEvalRes := domain.EvaluationResult{
		SubmissionId: sid,
		EvalStatus:   domain.PENDING,
		TL:           TL,
		ML:           ML,
		NCases:       nCase,
	}

	got, err := er.collection.InsertOne(ctx, newEvalRes)
	if err != nil {
		return "", nil
	}

	log := config.GetLogger()
	log.Info().Msgf("Created new eval result with id: [%s] to the database", got.InsertedID.(bson.ObjectID).Hex())

	var updated domain.EvaluationResult
	log.Debug().Msgf("Submisison bson: %v", sid)
	if err := er.collection.FindOne(ctx, bson.M{"submission_id": sid}).Decode(&updated); err != nil {
		return "", err
	}

	log.Debug().Msgf("Item: %v", updated)
	return got.InsertedID.(bson.ObjectID).Hex(), nil
}

func (er *EvaluationRepositoryImpl) UpdateVerdict(ctx context.Context, evalId string, vert domain.Verdict) error {
	bid, err := bson.ObjectIDFromHex(evalId)
	if err != nil {
		return err
	}

	_, err = er.collection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{"$set": bson.M{"verdict": vert}})
	if err != nil {
		return err
	}
	return nil
}

func (er *EvaluationRepositoryImpl) GetEval(ctx context.Context, evalId string) (*domain.EvaluationResult, error) {
	bid, err := bson.ObjectIDFromHex(evalId)
	if err != nil {
		return nil, err
	}
	var returnEval domain.EvaluationResult
	if err := er.collection.FindOne(ctx, bson.M{"_id": bid}).Decode(&returnEval); err != nil {
		return nil, err
	}

	return &returnEval, nil
}

func (er *EvaluationRepositoryImpl) GetEvalBson(ctx context.Context, bid bson.ObjectID) (*domain.EvaluationResult, error) {
	var returnEval domain.EvaluationResult
	if err := er.collection.FindOne(ctx, bson.M{"_id": bid}).Decode(&returnEval); err != nil {
		return nil, err
	}

	return &returnEval, nil
}

func (er *EvaluationRepositoryImpl) GetEvalBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.EvaluationResult, error) {
	var returnEval domain.EvaluationResult
	if err := er.collection.FindOne(ctx, bson.M{"submission_id": submissionId}).Decode(&returnEval); err != nil {
		return nil, err
	}

	return &returnEval, nil
}

func (er *EvaluationRepositoryImpl) GetEvalBySubmissionIdNoBson(ctx context.Context, submissionId string) (*domain.EvaluationResult, error) {
	submissionBsonId, err := bson.ObjectIDFromHex(submissionId)
	if err != nil {
		return nil, err
	}
	return er.GetEvalBySubmissionId(ctx, submissionBsonId)
}

func (er *EvaluationRepositoryImpl) UpdateCase(
	ctx context.Context,
	evalId string,
	verdictCase domain.Verdict,
	cpuTimeCase float64,
	memoryUsageCase memory.Memory,
	outputCase string,
	pointsCase int,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
) error {
	bid, err := bson.ObjectIDFromHex(evalId)
	if err != nil {
		return err
	}

	_, err = er.collection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{
		"$push": bson.M{
			"verdict_case":      verdictCase,
			"cpu_time_case":     cpuTimeCase,
			"memory_usage_case": memoryUsageCase,
			"outputs":           outputCase,
			"points_case":       pointsCase,
		},
		"$set": bson.M{
			"eval_status":  domain.JUDGING,
			"cpu_time":     cpuTime,
			"memory_usage": memoryUsage,
			"n_success":    nsucess,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (er *EvaluationRepositoryImpl) UpdateFinal(
	ctx context.Context,
	evalId string,
	verdict domain.Verdict,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
	points int,
	message string,
) error {
	bid, err := bson.ObjectIDFromHex(evalId)
	if err != nil {
		return err
	}

	_, err = er.collection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{"$set": bson.M{
		"verdict":          verdict,
		"cpu_time":         cpuTime,
		"memory_usage":     memoryUsage,
		"n_success":        nsucess,
		"points":           points,
		"message":          message,
		"eval_status":      domain.FINISHED,
		"timestamp_finish": time.Now().UnixMilli(),
	}})

	if err != nil {
		return err
	}
	return nil
}

func (er *EvaluationRepositoryImpl) UpdateCaseFloat(
	ctx context.Context,
	evalId string,
	verdictCase domain.Verdict,
	cpuTimeCase float64,
	memoryUsageCase memory.Memory,
	outputCase string,
	scoreCase float64,
	maxScoreCase float64,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
) error {
	bid, err := bson.ObjectIDFromHex(evalId)
	if err != nil {
		return err
	}

	_, err = er.collection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{
		"$push": bson.M{
			"verdict_case":      verdictCase,
			"cpu_time_case":     cpuTimeCase,
			"memory_usage_case": memoryUsageCase,
			"outputs":           outputCase,
			"max_score_case":    maxScoreCase,
			"score_case":        scoreCase,
		},
		"$set": bson.M{
			"eval_status":  domain.JUDGING,
			"cpu_time":     cpuTime,
			"memory_usage": memoryUsage,
			"n_success":    nsucess,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (er *EvaluationRepositoryImpl) UpdateFinalFloat(
	ctx context.Context,
	evalId string,
	verdict domain.Verdict,
	cpuTime float64,
	memoryUsage memory.Memory,
	nsucess int,
	totalScore float64,
	totalMaxScore float64,
	message string,
) error {
	bid, err := bson.ObjectIDFromHex(evalId)
	if err != nil {
		return err
	}

	_, err = er.collection.UpdateOne(ctx, bson.M{"_id": bid}, bson.M{"$set": bson.M{
		"verdict":          verdict,
		"cpu_time":         cpuTime,
		"memory_usage":     memoryUsage,
		"n_success":        nsucess,
		"score":            totalScore,
		"max_score":        totalMaxScore,
		"message":          message,
		"eval_status":      domain.FINISHED,
		"timestamp_finish": time.Now().UnixMilli(),
	}})

	if err != nil {
		return err
	}
	return nil
}

// FilterJudgedSubmissions filters submissions where it has been judged (FINISHED status)
// func (er *EvaluationRepositoryImpl) FilterJudgedSubmissions(ctx context.Context, submissionIds []string) ([]string, error) {
// 	var bsonIds []bson.ObjectID
// 	for _, id := range submissionIds {
// 		bid, err := bson.ObjectIDFromHex(id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		bsonIds = append(bsonIds, bid)
// 	}
//
// 	filter := bson.M{
// 		"$and": bson.A{
// 			bson.M{"eval_status": domain.FINISHED},
// 			bson.M{"submission_id": bson.M{"$in": bsonIds}},
// 		},
// 	}
// 	cursor, err := er.collection.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer cursor.Close(ctx)
//
// 	var filteredIds []string
// 	for cursor.Next(ctx) {
// 		var doc struct {
// 			submissionID bson.ObjectID `bson:"submission_id"`
// 		}
//
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, err
// 		}
// 		filteredIds = append(filteredIds, doc.submissionID.Hex())
// 	}
//
// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}
// 	return filteredIds, nil
// }
