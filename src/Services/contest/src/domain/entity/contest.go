package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	ProblemCountLimit uint8 = 15
)

const (
	ScoreboardHidden         string = "SCOREBOARD_HIDDEN"
	ScoreboardPublic         string = "SCOREBOARD_PUBLIC"
	ScoreboardContestantOnly string = "SCOREBOARD_CONTESTANT_ONLY"
)

type Contest struct {
	Id          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name" json:"name,omitempty"`
	Description string        `bson:"description" json:"description,omitempty"`

	Authors     []string     `bson:"authors" json:"authors,omitempty"`
	Curators    []string     `bson:"curators" json:"curators,omitempty"`
	Testers     []string     `bson:"testers" json:"testers,omitempty"`
	Contestants []Contestant `bson:"contestants" json:"contestants,omitempty"`

	ProblemLabels []string `bson:"problem_labels" json:"problem_labels,omitempty"`
	Problems      []uint64 `bson:"problems" json:"problems,omitempty"`

	ScoreboardVisibility string `bson:"scoreboard_visibility" json:"scoreboard_visibility,omitempty"`

	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`

	ScoringType ScoringType `bson:"scoring_type" json:"scoring_type,omitempty"`
	ICPCRules   ICPCRule    `bson:"penalty_rules" json:"penalty_rules"`
	IOIRules    IOIRule     `bson:"ioi_rules" json:"ioi_rules"`

	Status           ContestStatus `bson:"status" json:"status,omitempty"`
	FinalizeAt       time.Time     `bson:"finalize_at" json:"finalize_at"`
	RejudgeWindowEnd time.Time     `bson:"rejudge_window_end" json:"rejudge_window_end"`
}

type ICPCRule struct {
	PenaltyMinutes  uint16        `bson:"penalty_minutes" json:"penalty_minutes,omitempty"`
	FreezeStartTime time.Time     `bson:"freeze_start_time" json:"freeze_start_time"`
	FreezeTime      time.Duration `bson:"freeze_time" json:"freeze_time,omitempty"`
}

type IOIRule struct {
	MaxAllowedSubmissionsPerProblem uint16 `bson:"max_allowed_submissions_per_problem" json:"max_allowed_submissions_per_problem,omitempty"`
}

type ScoringType string

const (
	IOI  ScoringType = "IOI"
	ICPC ScoringType = "ICPC"
)

type ContestStatus string

const (
	ContestStatusDraft     ContestStatus = "DRAFT"
	ContestStatusScheduled ContestStatus = "SCHEDULED"
	ContestStatusRunning   ContestStatus = "RUNNING"
	ContestStatusFreeze    ContestStatus = "FREEZE"
	ContestStatusEnded     ContestStatus = "ENDED"
)

func (contest *Contest) ContestantExist(username string) bool {
	fmt.Printf("contestants: %v\n", contest.Contestants)
	for _, contestant := range contest.Contestants {
		if contestant.Username == username {
			return true
		}
	}
	return false
}

func (contest *Contest) clean() error {
	if contest.EndTime.IsZero() {
		return fmt.Errorf("invalid start time")
	}

	if contest.EndTime.IsZero() {
		return fmt.Errorf("invalid end time")
	}

	if !contest.StartTime.Before(contest.EndTime) {
		return fmt.Errorf("start time bigger than end time")
	}

	return nil
}

func (contest *Contest) hasStarted() bool {
	return time.Now().After(contest.StartTime)
}

func (contest *Contest) hasEnded() bool {
	return time.Now().After(contest.EndTime)
}

func (contest *Contest) canSeeScoreboard(username string) bool {
	if contest.ScoreboardVisibility == ScoreboardHidden || time.Now().Before(contest.StartTime) {
		return false
	}

	if contest.ScoreboardVisibility == ScoreboardPublic {
		return true
	}

	// if contest.ScoreboardVisibility == ScoreboardContestantOnly && slices.Contains(contest.Contestants, username) {
	// 	return true
	// }

	return false
}
