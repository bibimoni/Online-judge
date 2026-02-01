package domain

import "time"

type ContestRule struct {
	ScoringType ScoringType `bson:"scoring_type" json:"scoring_type,omitempty"`

	// ICPC Rules
	PenaltyMinutes  uint16        `bson:"penalty_minutes" json:"penalty_minutes,omitempty"`
	FreezeStartTime time.Time     `bson:"freeze_start_time" json:"freeze_start_time"`
	FreezeTime      time.Duration `bson:"freeze_time" json:"freeze_time,omitempty"`

	// IOI Rules
	MaxAllowedSubmissionsPerProblem uint16 `bson:"max_allowed_submissions_per_problem" json:"max_allowed_submissions_per_problem,omitempty"`
}

type ScoringType string

const (
	IOI  ScoringType = "IOI"
	ICPC ScoringType = "ICPC"
)
