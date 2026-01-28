package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Contestant struct {
	Username  string    `bson:"username"`
	RealStart time.Time `bson:"real_start"`

	Submissions []bson.ObjectID `bson:"submissions"`

	TotalPoints     float64            `bson:"total_points"`
	Points          map[uint64]float64 `bson:"points"` // points of each problem
	ParticipantType ParticipantType    `bson:"participant_type"`
}

type ParticipantType string

const (
	RatedParticipant   ParticipantType = "RATED"
	UnratedParticipant ParticipantType = "UNRATED"
	VirtualParticipant ParticipantType = "VIRTUAL"
)

var ParticipantTypes = []ParticipantType{RatedParticipant, UnratedParticipant, VirtualParticipant}

func CreateContestant(username string, participantType ParticipantType) Contestant {
	newContestant := Contestant{
		Username:  username,
		RealStart: time.Now(),

		Submissions: []bson.ObjectID{},

		TotalPoints:     0.00,
		Points:          map[uint64]float64{},
		ParticipantType: participantType,
	}
	return newContestant
}
