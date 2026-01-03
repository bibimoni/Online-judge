package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type ContestProblem struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProblemId uint64        `json:"problem_id,omitempty" bson:"problem_id"`
	Label     string        `json:"label,omitempty" bson:"label"`
	MaxPoints float64       `json:"max_points,omitempty" bson:"max_points"`
}
