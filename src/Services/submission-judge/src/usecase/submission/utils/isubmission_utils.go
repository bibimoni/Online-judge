package isubmission_utils

import (
	"context"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	erepository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/evaluation"
	screpository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/sourcecode"
	srepository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/submission"
	usecasews "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetSubmissionRequests converts submission IDs, filtered out valid submission (evalStatus is FINISHED)
// convert it into SubmissionRequest slice and return it
// TODO: move it to a repository layer
func GetSubmissionRequests(
	ctx context.Context,
	evalRepo erepository.EvaluationRepository,
	sourcecodeRepo screpository.SourcecodeRepository,
	submissionRepo srepository.SubmissionRepository,
	submissionIds []string,
) ([]isolateservice.SubmissionRequest, error) {
	bIds := make([]bson.ObjectID, 0, len(submissionIds))
	for _, id := range submissionIds {
		bId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		bIds = append(bIds, bId)
	}

	pipeline := mongo.Pipeline{
		// Match submissions with given IDs
		bson.D{{Key: "$match", Value: bson.M{"_id": bson.M{"$in": bIds}}}},
		// Lookup evaluations
		bson.D{
			{
				Key: "$lookup",
				Value: bson.M{
					"from":         evalRepo.GetCollectionName(),
					"localField":   "_id",
					"foreignField": "submission_id",
					"as":           "evaluations",
				},
			},
		},
		// Unwind evaluations array
		{{Key: "$unwind", Value: "$evaluations"}},
		// Match evaluations with status FINISHED
		{{Key: "$match", Value: bson.M{"evaluations.eval_status": "FINISHED"}}},
		// Find source codes for the submissions
		bson.D{
			{
				Key: "$lookup",
				Value: bson.M{
					"from":         sourcecodeRepo.GetCollectionName(),
					"localField":   "_id",
					"foreignField": "submission_id",
					"as":           "sourcecodes",
				},
			},
		},
		// Unwind sourcecodes array
		{{Key: "$unwind", Value: "$sourcecodes"}},
	}

	cursor, err := submissionRepo.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissionRequests []isolateservice.SubmissionRequest
	for cursor.Next(ctx) {
		var doc struct {
			ID          bson.ObjectID `bson:"_id"`
			ProblemId   string        `bson:"problem_id"`
			Username    string        `bson:"username"`
			Type        string        `bson:"type"`
			Timestamp   time.Time     `bson:"timestamp"`
			Sourcecodes struct {
				Language   string `bson:"language"`
				Sourcecode string `bson:"source_code"`
			} `bson:"sourcecodes"`
			Evaluations struct {
				EvaluationId bson.ObjectID `bson:"_id"`
				Verdict      string        `bson:"verdict"`
			} `bson:"evaluations"`
		}

		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		submissionRequests = append(submissionRequests, isolateservice.SubmissionRequest{
			SubmissionId:   doc.ID.Hex(),
			ProblemId:      doc.ProblemId,
			Username:       doc.Username,
			SubmissionType: domain.SubmissionType(doc.Type),
			LanguageId:     doc.Sourcecodes.Language,
			Sourcecode:     doc.Sourcecodes.Sourcecode,
			EvalId:         doc.Evaluations.EvaluationId.Hex(),
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return submissionRequests, nil
}

func GetSubmission(
	ctx context.Context,
	evalRepo erepository.EvaluationRepository,
	sourcecodeRepo screpository.SourcecodeRepository,
	submissionRepo srepository.SubmissionRepository,
	submissionId string,
) (*usecase.GetSubmissionOutput, error) {
	log := config.GetLogger()
	sub, err := submissionRepo.FindSubmission(ctx, submissionId)
	if err != nil {
		log.Debug().Msgf("Find submission, error: %v", err)
		return nil, err
	}

	source, err := sourcecodeRepo.GetSourceBySubmissionId(ctx, (*sub).Id)
	if err != nil {
		log.Debug().Msgf("Find sourcecode, error: %v", err)
		return nil, err
	}

	eval, err := evalRepo.GetEvalBySubmissionId(ctx, (*sub).Id)
	if err != nil {
		log.Debug().Msgf("Find eval, error: %v", err)
		return nil, err
	}

	lang, err := store.DefaultStore.Get((*source).LanguageId)
	if err != nil {
		return nil, err
	}

	returnVal := usecase.GetSubmissionOutput{
		ProblemId:       (*sub).ProblemId,
		Verdict:         (*eval).Verdict,
		VerdictCase:     (*eval).VerdictCase,
		CpuTime:         (*eval).CpuTime,
		CpuTimeCase:     (*eval).CpuTimeCase,
		MemoryUsage:     (*eval).MemoryUsage,
		MemoryUsageCase: (*eval).MemoryUsageCase,
		NSuccess:        (*eval).NSuccess,
		Outputs:         (*eval).Outputs,
		Message:         (*eval).Message,
		Points:          (*eval).Points,
		PointsCase:      (*eval).PointsCase,
		NCases:          (*eval).NCases,
		TL:              (*eval).TL,
		ML:              (*eval).ML,
		Username:        (*sub).Username,
		Timestamp:       (*sub).Timestamp,
		Type:            (*sub).Type,
		Language:        lang.DisplayName(),
		SourceCode:      (*source).SourceCode,
		EvalStatus:      (*eval).EvalStatus,
	}

	return &returnVal, nil
}

func GetSubmissionWithoutSourceCode(
	ctx context.Context,
	evalRepo erepository.EvaluationRepository,
	sourcecodeRepo screpository.SourcecodeRepository,
	submissionRepo srepository.SubmissionRepository,
	submissionId string,
) (*usecasews.WSSubmissionResponse, error) {

	sub, err := GetSubmission(
		ctx,
		evalRepo,
		sourcecodeRepo,
		submissionRepo,
		submissionId,
	)

	if err != nil {
		return nil, err
	}
	returnVal := usecasews.WSSubmissionResponse{
		Username:        sub.Username,
		SubmissionId:    submissionId,
		ProblemId:       sub.ProblemId,
		Timestamp:       sub.Timestamp,
		Language:        sub.Language,
		Verdict:         sub.Verdict,
		VerdictCase:     sub.VerdictCase,
		CpuTime:         sub.CpuTime,
		CpuTimeCase:     sub.CpuTimeCase,
		MemoryUsage:     sub.MemoryUsage,
		MemoryUsageCase: sub.MemoryUsageCase,
		NSuccess:        sub.NSuccess,
		Outputs:         sub.Outputs,
		Points:          sub.Points,
		PointsCase:      sub.PointsCase,
		Message:         sub.Message,
		EvalStatus:      sub.EvalStatus,
	}

	return &returnVal, nil
}
