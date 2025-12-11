package domain

import (
	"fmt"
	"time"
)

const (
	PROBLEM_COUNT_LIMIT uint8 = 15
)

const (
	SCOREBOARD_HIDDEN          string = "SCOREBOARD_HIDDEN"
	SCOREBOARD_PUBLIC          string = "SCOREBOARD_PUBLIC"
	SCOREBOARD_CONTESTANT_ONLY string = "SCOREBOARD_CONTESTANT_ONLY"
)

type Contest struct {
	Id          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`

	Authors     []uint64     `bson:"authors"`
	Curators    []uint64     `bson:"curators"`
	Testers     []uint64     `bson:"testers"`
	Contestants []Contestant `bson:"contestants"`

	ProblemLabels []string `bson:"problem_labels"`
	Problems      []uint64 `bson:"problems"`

	ScoreboardVisibility string `bson:"scoreboard_visibility"`

	StartTime time.Time `bson:"start_time"`
	EndTime   time.Time `bson:"end_time"`
}

func (contest *Contest) ContestantExist(userId uint64) bool {
	fmt.Printf("contestants: %v\n", contest.Contestants)
	for _, contestant := range contest.Contestants {
		if contestant.UserID == userId {
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

func (contest *Contest) has_started() bool {
	return time.Now().After(contest.StartTime)
}

func (contest *Contest) has_ended() bool {
	return time.Now().After(contest.EndTime)
}

func (contest *Contest) can_see_scoreboard(userId uint64) bool {
	if contest.ScoreboardVisibility == SCOREBOARD_HIDDEN || time.Now().Before(contest.StartTime) {
		return false
	}

	if contest.ScoreboardVisibility == SCOREBOARD_PUBLIC {
		return true
	}

	// if contest.ScoreboardVisibility == SCOREBOARD_CONTESTANT_ONLY && slices.Contains(contest.Contestants, UserId) {
	// 	return true
	// }

	return false
}
