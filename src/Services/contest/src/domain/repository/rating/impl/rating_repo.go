package ratingrepoimpl

import (
	domain "contest/src/domain/entity"
	ratingrepo "contest/src/domain/repository/rating"
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RatingRepositoryImpl struct {
	ratingResultCol  *mongo.Collection
	userRatingCol    *mongo.Collection
}

func NewRatingRepository(db *mongo.Database) ratingrepo.RatingRepository {
	return &RatingRepositoryImpl{
		ratingResultCol:  db.Collection("RatingResult"),
		userRatingCol:    db.Collection("UserRatingState"),
	}
}

// ---------- RatingResult operations ----------

func (r *RatingRepositoryImpl) SaveRatingResults(ctx context.Context, results []domain.RatingResult) error {
	if len(results) == 0 {
		return nil
	}

	docs := make([]any, len(results))
	for i := range results {
		results[i].CreatedAt = time.Now()
		docs[i] = results[i]
	}

	_, err := r.ratingResultCol.InsertMany(ctx, docs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to batch-insert rating results")
		return err
	}

	log.Info().Msgf("Saved %d rating results for contest %s", len(results), results[0].ContestId.Hex())
	return nil
}

func (r *RatingRepositoryImpl) GetRatingResultsByContest(ctx context.Context, contestId string) ([]domain.RatingResult, error) {
	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return nil, ratingrepo.ErrInvalidContestId
	}

	filter := bson.M{"contest_id": cId}
	opts := options.Find().SetSort(bson.D{{Key: "rank", Value: 1}})

	cursor, err := r.ratingResultCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.RatingResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *RatingRepositoryImpl) GetRatingHistory(ctx context.Context, username string) ([]domain.RatingResult, error) {
	filter := bson.M{"username": username}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.ratingResultCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.RatingResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// ---------- UserRatingState operations ----------

func (r *RatingRepositoryImpl) SaveUserRatingStates(ctx context.Context, states []domain.UserRatingState) error {
	if len(states) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, len(states))
	for i, st := range states {
		filter := bson.M{"username": st.Username}
		update := bson.M{"$set": st}
		models[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	}

	_, err := r.userRatingCol.BulkWrite(ctx, models)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bulk-upsert user rating states")
		return err
	}

	log.Info().Msgf("Upserted %d user rating states", len(states))
	return nil
}

func (r *RatingRepositoryImpl) GetUserRatingState(ctx context.Context, username string) (*domain.UserRatingState, error) {
	filter := bson.M{"username": username}

	var state domain.UserRatingState
	err := r.userRatingCol.FindOne(ctx, filter).Decode(&state)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ratingrepo.ErrUserRatingNotFound
		}
		return nil, err
	}
	return &state, nil
}

func (r *RatingRepositoryImpl) BatchGetUserRatingStates(ctx context.Context, usernames []string) (map[string]domain.UserRatingState, error) {
	if len(usernames) == 0 {
		return map[string]domain.UserRatingState{}, nil
	}

	filter := bson.M{"username": bson.M{"$in": usernames}}
	cursor, err := r.userRatingCol.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]domain.UserRatingState, len(usernames))
	for cursor.Next(ctx) {
		var state domain.UserRatingState
		if err := cursor.Decode(&state); err != nil {
			return nil, err
		}
		result[state.Username] = state
	}

	return result, cursor.Err()
}
