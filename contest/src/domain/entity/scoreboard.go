package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ScoreboardSnapshot struct {
	Id        bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	ContestId bson.ObjectID   `json:"contest_id" bson:"contest_id"`
	Kind      SnapshotKind    `json:"kind" bson:"kind"`
	Rows      []ScoreboardRow `json:"rows" bson:"rows"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
}

type SnapshotKind string

const (
	LiveSnapshot  SnapshotKind = "LIVE"
	FinalSnapshot SnapshotKind = "FINAL"
)

type ScoreboardRow struct {
	Username string                    `json:"username" bson:"username"`
	Rank     int                       `json:"rank" bson:"rank"`
	Score    float64                   `json:"score" bson:"score"`
	Penalty  int                       `json:"penalty,omitempty" bson:"penalty,omitempty"` // In minutes for ICPC
	Problems []ScoreboardProblemResult `json:"problems" bson:"problems"`
}

type ScoreboardProblemResult struct {
	Label            string  `json:"label" bson:"label"`                         // Problem label (A, B, C, etc.)
	ProblemId        uint64  `json:"problem_id" bson:"problem_id"`               // Reference to problem
	Points           float64 `json:"points" bson:"points"`                       // Points earned (IOI) or 1/0 (ICPC)
	Attempts         int     `json:"attempts" bson:"attempts"`                   // Number of submissions
	Solved           bool    `json:"solved" bson:"solved"`                       // True if accepted
	SolvedAt         int     `json:"solved_at,omitempty" bson:"solved_at,omitempty"`       // Minutes from contest start
	Pending          bool    `json:"pending,omitempty" bson:"pending,omitempty"`                     // Has pending/judging submissions
	FirstSolve       bool    `json:"first_solve,omitempty" bson:"first_solve,omitempty"`   // First to solve (optional)
	PenaltyAttempts  int     `json:"penalty_attempts,omitempty" bson:"penalty_attempts,omitempty"`   // Failed attempts before AC (ICPC)
	BestContestSubmissionId bson.ObjectID  `json:"best_contest_submission_id,omitempty" bson:"best_contest_submission_id,omitempty"` // For IOI
}

type ICPCScoreboardRow struct {
	ScoreboardRow
	TotalSolved int `json:"total_solved" bson:"total_solved"` 
}

type IOIScoreboardRow struct {
	ScoreboardRow
	TotalPoints float64 `json:"total_points" bson:"total_points"` 
}
