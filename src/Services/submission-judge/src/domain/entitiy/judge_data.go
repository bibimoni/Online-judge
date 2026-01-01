package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type JudgeData struct {
	Id        bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	ProblemId bson.ObjectID   `json:"problem_id,omitempty" bson:"problem_id"`
	Format    JudgeDataFormat `json:"format,omitempty" bson:"format"`
}

type JudgeDataFormat string
