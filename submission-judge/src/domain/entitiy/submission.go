package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Submission struct {
	Id        bson.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string         `json:"username,omitempty" bson:"username"`
	ProblemId string         `json:"problem_id,omitempty" bson:"problem_id"`
	Timestamp time.Time      `json:"timestamp" bson:"timestamp"`
	Type      SubmissionType `json:"type,omitempty" bson:"type"`
	// JudgeDataId  *bson.ObjectID `json:"judge_data_id,omitempty" bson:"judge_data_id"`
	// MetricId     *bson.ObjectID `json:"metric_id,omitempty" bson:"metric_id"`
}

type SubmissionType string

const (
	CUSTOM SubmissionType = "CUSTOM"
	ICPC   SubmissionType = "ICPC"
	IOI    SubmissionType = "IOI"
)
