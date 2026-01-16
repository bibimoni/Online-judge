package domain

import (
	"contest/src/common"
	"fmt"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	ProblemCountLimit uint8 = 15
)

type ScoreboardVisibility string

const (
	ScoreboardHidden         ScoreboardVisibility = "SCOREBOARD_HIDDEN"
	ScoreboardPublic         ScoreboardVisibility = "SCOREBOARD_PUBLIC"
	ScoreboardContestantOnly ScoreboardVisibility = "SCOREBOARD_CONTESTANT_ONLY"
)

type Contest struct {
	Id          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name" json:"name,omitempty"`
	Description string        `bson:"description" json:"description,omitempty"`

	Authors     []string     `bson:"authors" json:"authors,omitempty"`
	Admins      []string     `bson:"admins" json:"curators,omitempty"`
	Testers     []string     `bson:"testers" json:"testers,omitempty"`
	Contestants []Contestant `bson:"contestants" json:"contestants,omitempty"`

	Problems []ContestProblem `bson:"problems" json:"problems,omitempty"`

	ScoreboardVisibility ScoreboardVisibility `bson:"scoreboard_visibility" json:"scoreboard_visibility,omitempty"`

	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`

	// ScoringType ScoringType `bson:"scoring_type" json:"scoring_type,omitempty"`

	// ICPCRules   ICPCRule    `bson:"penalty_rules" json:"penalty_rules"`
	// IOIRules    IOIRule     `bson:"ioi_rules" json:"ioi_rules"`

	ContestRule ContestRule `bson:"contest_rule" json:"contest_rule"`

	Status           ContestStatus     `bson:"status" json:"status,omitempty"`
	FinalizeAt       time.Time         `bson:"finalize_at" json:"finalize_at"`
	RejudgeWindowEnd time.Time         `bson:"rejudge_window_end" json:"rejudge_window_end"`
	Visibility       ContestVisibility `bson:"visibility" json:"visibility,omitempty"`
}

type ContestVisibility string

const (
	VisibilityPrivate ContestVisibility = "HIDDEN"
	VisibilityPublic  ContestVisibility = "PUBLIC"
)

type ContestStatus string

const (
	Draft     ContestStatus = "DRAFT"
	Scheduled ContestStatus = "SCHEDULED"
	Running   ContestStatus = "RUNNING"
	Freeze    ContestStatus = "FREEZE"
	Ended     ContestStatus = "ENDED"
)

func (contest *Contest) IsAdmin(username, role string) bool {
	return role == common.AdminRole || slices.Contains(contest.Admins, username)
}

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

func (contest *Contest) CanViewContest(username, role string) bool {
	if contest.Visibility == VisibilityPublic {
		return true
	}

	if slices.Contains(contest.Admins, username) ||
		slices.Contains(contest.Testers, username) ||
		slices.Contains(contest.Authors, username) ||
		contest.ContestantExist(username) ||
		role == common.AdminRole {
		return true
	}
	return false
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
