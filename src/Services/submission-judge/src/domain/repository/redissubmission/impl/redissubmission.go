package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	repository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/redissubmission"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
	"github.com/redis/go-redis/v9"
)

type RedisSubmissionRepositoryImpl struct {
	rdb *redis.Client
}

func NewRedisSubmissionRepositoryImpl(rdb *redis.Client) *RedisSubmissionRepositoryImpl {
	return &RedisSubmissionRepositoryImpl{
		rdb,
	}
}

func (rs *RedisSubmissionRepositoryImpl) GetChannelString(problemId, username, submissionId string) string {
	return fmt.Sprintf("%s:%s:%s", problemId, username, submissionId)
}

func NewRedisSubmissionRepository(rdb *redis.Client) repository.RedisSubmissionRepository {
	return NewRedisSubmissionRepositoryImpl(rdb)
}

func (rs *RedisSubmissionRepositoryImpl) PulishSubmission(ctx context.Context, res usecase.WSSubmissionResponse) error {
	channel := rs.GetChannelString(res.ProblemId, res.Username, res.SubmissionId)

	byte, err := json.Marshal(res)
	if err != nil {
		return err
	}
	config.GetLogger().Debug().Msgf("Channel: %s --- Receive event %v", channel, res)
	return rs.rdb.Publish(ctx, channel, byte).Err()
}

func (rs *RedisSubmissionRepositoryImpl) Subscribe(ctx context.Context, channelId string) (<-chan *usecase.WSSubmissionResponse, error) {
	sub := rs.rdb.PSubscribe(ctx, channelId)
	raw := sub.Channel()

	out := make(chan *usecase.WSSubmissionResponse)

	go func() {
		defer sub.Close()
		defer close(out)

		for {
			select {
			case msg, ok := <-raw:
				if !ok {
					return
				}
				var upd usecase.WSSubmissionResponse
				err := json.Unmarshal([]byte(msg.Payload), &upd)
				if err != nil {
					config.GetLogger().Error().Msgf("Json Unmarshal failed: %v", err)
					continue
				}

				out <- &upd
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func (rs *RedisSubmissionRepositoryImpl) PushSubmissionJob(ctx context.Context, req *isolateservice.SubmissionRequest) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	return pushJob(ctx, rs.rdb, cfg.Redis.SubmissionQueueKey, req)
}

func (rs *RedisSubmissionRepositoryImpl) PopSubmissionJob(ctx context.Context) (*isolateservice.SubmissionRequest, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return popJob[isolateservice.SubmissionRequest](ctx, rs.rdb, cfg.Redis.SubmissionQueueKey)
}

func (rs *RedisSubmissionRepositoryImpl) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return rs.rdb.SetNX(ctx, key, value, ttl).Result()
}

func pushJob[T any](ctx context.Context, rdb *redis.Client, queueKey string, req *T) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return rdb.RPush(ctx, queueKey, data).Err()
}

func popJob[T any](ctx context.Context, rdb *redis.Client, queueKey string) (*T, error) {
	res, err := rdb.BLPop(ctx, 0, queueKey).Result()
	if err != nil {
		return nil, err
	}

	var req T
	// res[1] is value
	err = json.Unmarshal([]byte(res[1]), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}
