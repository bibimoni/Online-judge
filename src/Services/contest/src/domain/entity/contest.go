package domain

import (
	"contest/src/common"
	"contest/src/infrastructure/config"
	"fmt"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	ProblemCountLimit uint8 = 15
	RegisterKey string = "register"
	UnRegisterKey string = "unregister"
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

func (contest *Contest) CanRegister(username string, _ ParticipantType) bool {
	// TODO: support participant type check, specificlly for virtual contest
	if contest.ContestantExist(username) {
		return false;
	}

	canRegisterPhase := []ContestStatus{Scheduled, Running}
	if contest.IsContestManager(username) {
		canRegisterPhase = append(canRegisterPhase, Draft)
	}

	if !slices.Contains(canRegisterPhase, contest.Status) {
		return false
	}

	return true
}

func (contest *Contest) CanUnregister(username string) bool {
	if !contest.ContestantExist(username) {
		return false;
	}
	canUnregisterPhase := []ContestStatus{Scheduled}
	if contest.IsContestManager(username) {
		canUnregisterPhase = append(canUnregisterPhase, Draft)
	}

	if !slices.Contains(canUnregisterPhase, contest.Status) {
		return false
	}

	return true
}

func (contest *Contest) CanSubmit(username string) bool {
	if !contest.ContestantExist(username) {
		return false
	}
	canSubmitPhase := []ContestStatus{Running, Freeze}
	if contest.IsContestManager(username) {
		canSubmitPhase = append(canSubmitPhase, Draft, Scheduled)
	}

	if !slices.Contains(canSubmitPhase, contest.Status) {
		return false
	}

	return true
}

func (contest *Contest) IsAdmin(username, role string) bool {
	return role == common.AdminRole || slices.Contains(contest.Admins, username)
}

func (contest *Contest) ContestantExist(username string) bool {
	config.GetLogger().Debug().Msgf("contestants: %v\n", contest.Contestants)
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

func (contest *Contest) IsContestManager(username string) bool {
	return slices.Contains(contest.Admins, username) ||
		slices.Contains(contest.Testers, username) ||
		slices.Contains(contest.Authors, username)
}

func (contest *Contest) CanViewContest(username, role string) bool {
	if contest.Visibility == VisibilityPublic {
		return true
	}

	if contest.IsContestManager(username) ||
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
