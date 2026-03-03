package scoreboardusecase

import (
	domain "contest/src/domain/entity"
	"context"
	"time"
)

type ScoreboardInteractor interface {
	BuildScoreboardSnapshot(ctx context.Context, input *BuildScoreboardSnapshotInput) (*BuildScoreboardSnapshotOutput, error)
	GetScoreboardAt(ctx context.Context, input *GetScoreboardAtInput) (*BuildScoreboardSnapshotOutput, error)
}

type (
	BuildScoreboardSnapshotInput struct {
		ContestId      string `json:"contest_id"`
		IncludeVirtual bool	  `json:"include_virtual"`
		IncludeUnrated bool	  `json:"include_unrated"`
		Username       string 
		UserRole       string 
		Authenticated  bool
	}
	BuildScoreboardSnapshotOutput struct {
		ScoreboardSnapshot *domain.ScoreboardSnapshot `json:"scoreboard_snapshot"`
	}
	GetScoreboardAtInput struct {
		ContestId      string    `json:"contest_id"`
		At             time.Time `json:"at"`
		IncludeVirtual bool      `json:"include_virtual"`
		IncludeUnrated bool      `json:"include_unrated"`
		Username       string
		UserRole       string
		Authenticated  bool
	}
)
