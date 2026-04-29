package scoreboardusecase

import (
	domain "contest/src/domain/entity"
	"context"
)

type ScoreboardInteractor interface {
	BuildScoreboardSnapshot(ctx context.Context, input *BuildScoreboardSnapshotInput) (*BuildScoreboardSnapshotOutput, error)
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
)
