package scoreboardrepo

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

type ScoreboardRepository interface {
	SaveSnapshot(ctx context.Context, snapshot *domain.ScoreboardSnapshot) error
	GetSnapshot(ctx context.Context, contestId string, kind domain.SnapshotKind) (*domain.ScoreboardSnapshot, error)
	GetLatestSnapshot(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error)
	SaveLiveScoreboard(ctx context.Context, contestId string, snapshot *domain.ScoreboardSnapshot) error
	GetLiveScoreboard(ctx context.Context, contestId string) (*domain.ScoreboardSnapshot, error)
	PublishScoreboardUpdate(ctx context.Context, contestId string, version int64) error
	DeleteLiveScoreboard(ctx context.Context, contestId string) error
	CreateAndGetSnapshot(
		ctx context.Context, 
		contestId string,
		kind domain.SnapshotKind,
		rows []domain.ScoreboardRow,
	) (*domain.ScoreboardSnapshot, error)
}

var (
	ErrSnapshotNotFound      = errors.New("scoreboard snapshot not found")
	ErrInvalidContestId      = errors.New("invalid contest id")
	ErrScoreboardUnavailable = errors.New("scoreboard unavailable")
)
