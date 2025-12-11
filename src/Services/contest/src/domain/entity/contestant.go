package domain

import "time"

type Contestant struct {
	UserID    uint64    `bson:"user-id"`
	RealStart time.Time `bson:"real-start"`

	Submissions []uint64 `bson:"submissions"`

	TotalPoints float64            `bson:"total-points"`
	Points      map[uint64]float64 `bson:"points"` // points of each problem
}

func CreateContestant(userId uint64) Contestant {
	newContestant := Contestant{
		UserID:    userId,
		RealStart: time.Now(),

		Submissions: []uint64{},

		TotalPoints: 0.00,
		Points:      map[uint64]float64{},
	}
	return newContestant
}
