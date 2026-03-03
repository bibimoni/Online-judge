package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// RatingResult stores the per-user rating change produced by a single contest.
type RatingResult struct {
	Id          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ContestId   bson.ObjectID `json:"contest_id" bson:"contest_id"`
	Username    string        `json:"username" bson:"username"`
	Rank        int           `json:"rank" bson:"rank"`
	OldRating   float64       `json:"old_rating" bson:"old_rating"`     // display rating before
	NewRating   float64       `json:"new_rating" bson:"new_rating"`     // display rating after
	Delta       float64       `json:"delta" bson:"delta"`               // NewRating - OldRating
	Performance float64       `json:"performance" bson:"performance"`   // single-round performance
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}

// ---------- Elo-MMR specific types ----------

// EloMMRConfig holds the tunable hyper-parameters for Elo-MMR.
type EloMMRConfig struct {
	// Mu0 is the initial mean rating for a new player (default 1500).
	Mu0 float64
	// Sig0 is the initial uncertainty (standard deviation) for a new player (default 350).
	Sig0 float64
	// SigLimit is a lower bound on sigma to prevent overconfidence (default 80).
	SigLimit float64
	// DriftPerSec is the variance growth rate per second of inactivity,
	// modelling skill drift over time.  A reasonable default is 0 (disabled)
	// or a small value like 1e-5.
	DriftPerSec float64
	// Beta is the performance-noise parameter.  200 is a reasonable competitive-
	// programming default (similar to Codeforces / AtCoder spread).
	Beta float64
	// WeightLimit caps the total logistic weight assigned to a single contest
	// (used in the robust weight computation).  MaxContests equivalent.
	// A good default is 6.0.
	WeightLimit float64
}

// DefaultEloMMRConfig returns production-ready defaults.
func DefaultEloMMRConfig() EloMMRConfig {
	return EloMMRConfig{
		Mu0:         1500.0,
		Sig0:        350.0,
		SigLimit:    80.0,
		DriftPerSec: 0.0,
		Beta:        200.0,
		WeightLimit: 6.0,
	}
}

// UserRatingState is the persistent per-player state that the Elo-MMR algorithm
// updates after every contest.
type UserRatingState struct {
	Username       string    `json:"username" bson:"username"`
	Mu             float64   `json:"mu" bson:"mu"`                         // current Gaussian mean
	Sigma          float64   `json:"sigma" bson:"sigma"`                   // current Gaussian std-dev
	DisplayRating  float64   `json:"display_rating" bson:"display_rating"` // conservative estimate shown to user
	NumContests    int       `json:"num_contests" bson:"num_contests"`     // number of rated contests
	LastContestAt  time.Time `json:"last_contest_at" bson:"last_contest_at"`
}

// NewUserRatingState creates an initial state for a first-time participant.
func NewUserRatingState(username string, cfg EloMMRConfig) UserRatingState {
	return UserRatingState{
		Username:      username,
		Mu:            cfg.Mu0,
		Sigma:         cfg.Sig0,
		DisplayRating: cfg.Mu0 - 2*cfg.Sig0, // conservative lower bound
		NumContests:   0,
	}
}

// StandingRow is the minimal per-user input the algorithm needs from a contest scoreboard.
type StandingRow struct {
	Username string
	Rank     int     // 1-based rank (ties share the same rank)
	Score    float64 // only used for tie-breaking display; rank is authoritative
}
