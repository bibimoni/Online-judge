package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Problem struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ProblemId     uint64             `json:"problem-id" bson:"problem-id"`
	Name          string             `json:"name" bson:"name"`
	ShortName     string             `json:"short-name" bson:"short-name"`
	Tags          []string           `json:"tags" bson:"tags"`
	TestNum       uint64             `json:"test-num" bson:"test-num"`
	TimeLimit     uint64             `json:"time-limit" bson:"time-limit"`
	MemoryLimit   uint64             `json:"memory-limit" bson:"memory-limit"`
	IsInteractive bool               `json:"is-interactive" bson:"is-interactive"`

	ScoringMode   string         `json:"scoring-mode,omitempty"`
	TestGroups    []ProblemGroup `json:"test-groups,omitempty"`
	TestMaxScores []float64      `json:"test-max-scores,omitempty"`
}

type ProblemGroup struct {
	Name         string   `json:"name,omitempty"`
	Scoring      string   `json:"scoring,omitempty"`
	MaxScore     float64  `json:"max-score,omitempty"`
	TestIndices  []int    `json:"test-indices,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type Package struct {
	ID                  uint64 `json:"id"`
	Revision            uint64 `json:"revision"`
	CreationTimeSeconds uint64 `json:"creationTimeSeconds"`
	State               string `json:"state"`
	Comment             string `json:"comment"`
	Type                string `json:"type"`
}
