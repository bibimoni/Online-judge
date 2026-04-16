package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type RejudgeJob struct {
	Id            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ContestId     bson.ObjectID `json:"contest_id,omitempty" bson:"contest_id"`
	SubmissionIds []string      `json:"submission_ids,omitempty" bson:"submission_ids"`
	ProblemId     uint64        `json:"problem_id,omitempty" bson:"problem_id"`
	RequestedBy   string        `json:"requested_by,omitempty" bson:"requested_by"`
}
