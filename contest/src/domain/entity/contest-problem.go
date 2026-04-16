package domain

type ContestProblem struct {
	ProblemId uint64  `json:"problem_id,omitempty" bson:"problem_id"`
	Label     string  `json:"label,omitempty" bson:"label"`
	MaxPoints float64 `json:"max_points,omitempty" bson:"max_points"`
}
