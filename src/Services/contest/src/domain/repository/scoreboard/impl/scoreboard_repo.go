package scoreboardrepoimpl

import (
	domain "contest/src/domain/entity"
	scoreboardrepo "contest/src/domain/repository/scoreboard"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ScoreboardRepositoryImpl struct {
	collection *mongo.Collection
	redis      *redis.Client
}

func NewScoreboardRepository(db *mongo.Database, redis *redis.Client) scoreboardrepo.ScoreboardRepository {
	return &ScoreboardRepositoryImpl{
		collection: db.Collection("ScoreboardSnapshot"),
		redis:      redis,
	}
}

func (r *ScoreboardRepositoryImpl) CreateAndGetSnapshot(
	ctx context.Context, 
	contestId string,
	kind domain.SnapshotKind,
	rows []domain.ScoreboardRow,
) (*domain.ScoreboardSnapshot, error) {
	contestBsonId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return nil, scoreboardrepo.ErrInvalidContestId
	}
	newSnapshot := &domain.ScoreboardSnapshot{
		ContestId: contestBsonId,
		Kind:      kind,
		Rows:      rows,
		CreatedAt: time.Now(),
	}
	result, err := r.collection.InsertOne(ctx, newSnapshot)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create scoreboard snapshot for contest %s", contestId)
		return nil, err
	}
	insertedId, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		log.Error().Msg("Failed to convert inserted ID to bson.ObjectID")
		return nil, errors.New("failed to convert inserted ID to bson.ObjectID")
	}

	newSnapshot.Id = insertedId
	return newSnapshot, nil
}

func (r *ScoreboardRepositoryImpl) SaveSnapshot(ctx context.Context, snapshot *domain.ScoreboardSnapshot) error {
	if snapshot.ContestId.IsZero() {
		return scoreboardrepo.ErrInvalidContestId
	}

	snapshot.CreatedAt = time.Now()
	filter := bson.M{
		"contest_id": snapshot.ContestId,
		"kind":       snapshot.Kind,
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, filter, snapshot, opts)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save scoreboard snapshot for contest %s", snapshot.ContestId.Hex())
		return err
	}

	log.Info().Msgf("Saved %s scoreboard snapshot for contest %s", snapshot.Kind, snapshot.ContestId.Hex())
	return nil
}

func (r *ScoreboardRepositoryImpl) GetSnapshot(ctx context.Context, contestId string, kind domain.SnapshotKind) (*domain.ScoreboardSnapshot, error) {
	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return nil, scoreboardrepo.ErrInvalidContestId
	}

	filter := bson.M{
		"contest_id": cId,
		"kind":       kind,
	}

	var snapshot domain.ScoreboardSnapshot
	err = r.collection.FindOne(ctx, filter).Decode(&snapshot)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, scoreboardrepo.ErrSnapshotNotFound
		}
		return nil, err
	}

	return &snapshot, nil
}
func (r *ScoreboardRepositoryImpl) GetLatestSnapshot(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error) {
	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return nil, scoreboardrepo.ErrInvalidContestId
	}

	filter := bson.M{"contest_id": cId}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	var snapshot domain.ScoreboardSnapshot
	err = r.collection.FindOne(ctx, filter, opts).Decode(&snapshot)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, scoreboardrepo.ErrSnapshotNotFound
		}
		return nil, err
	}

	return &snapshot, nil
}
func (r *ScoreboardRepositoryImpl) SaveLiveScoreboard(ctx context.Context, contestId string, snapshot *domain.ScoreboardSnapshot) error {
	key := fmt.Sprintf("contest:scoreboard:live:%s", contestId)

	data, err := json.Marshal(snapshot)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal scoreboard snapshot")
		return err
	}
	err = r.redis.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save live scoreboard to Redis for contest %s", contestId)
		return err
	}

	log.Debug().Msgf("Saved live scoreboard to Redis for contest %s", contestId)
	return nil
}
func (r *ScoreboardRepositoryImpl) GetLiveScoreboard(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error) {
	key := fmt.Sprintf("contest:scoreboard:live:%s", contestId)

	data, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, scoreboardrepo.ErrSnapshotNotFound
		}
		log.Error().Err(err).Msgf("Failed to get live scoreboard from Redis for contest %s", contestId)
		return nil, err
	}

	var snapshot domain.ScoreboardSnapshot
	err = json.Unmarshal([]byte(data), &snapshot)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal scoreboard snapshot from Redis")
		return nil, err
	}

	return &snapshot, nil
}
func (r *ScoreboardRepositoryImpl) PublishScoreboardUpdate(ctx context.Context, contestId string, version int64) error {
	channel := fmt.Sprintf("contest:scoreboard:%s", contestId)

	message := map[string]interface{}{
		"contest_id":   contestId,
		"version":      version,
		"generated_at": time.Now().Unix(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.redis.Publish(ctx, channel, data).Err()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to publish scoreboard update for contest %s", contestId)
		return err
	}

	log.Debug().Msgf("Published scoreboard update event for contest %s (version %d)", contestId, version)
	return nil
}
func (r *ScoreboardRepositoryImpl) DeleteLiveScoreboard(ctx context.Context, contestId string) error {
	key := fmt.Sprintf("contest:scoreboard:live:%s", contestId)

	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete live scoreboard from Redis for contest %s", contestId)
		return err
	}

	log.Info().Msgf("Deleted live scoreboard from Redis for contest %s", contestId)
	return nil
}
