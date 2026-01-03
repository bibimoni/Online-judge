package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ContestSubmission struct {
	Id               bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ContestId        bson.ObjectID `json:"contest_id,omitempty" bson:"contest_id"`
	SubmissionId     bson.ObjectID `json:"submission_id,omitempty" bson:"submission_id"`
	Username         string        `json:"username,omitempty" bson:"username"`
	ContestProblemId bson.ObjectID `json:"contest_problem_id,omitempty" bson:"contest_problem_id"`
	SubmittedAt      time.Time     `json:"submitted_at" bson:"submitted_at"`
	Verdict          Verdict       `json:"verdict,omitempty" bson:"verdict"`
	Points           float64       `json:"points,omitempty" bson:"points"`
	JudgedAt         time.Time     `json:"judged_at" bson:"judged_at"`
	EvalStatus       EvalStatus    `json:"eval_status,omitempty" bson:"eval_status"`
	Ignored          bool          `json:"ignored,omitempty" bson:"ignored"`
	UpdatedAt        time.Time     `json:"updated_at" bson:"updated_at"`
}

type EvalStatus string

const (
	Pending  EvalStatus = "PENDING"
	Judging  EvalStatus = "JUDGING"
	Finished EvalStatus = "FINISHED"
)

type Verdict string

const (
	Accepted            Verdict = "ACCEPTED"
	CompilationError    Verdict = "COMPILATION_ERROR"
	Rejected            Verdict = "REJECTED"
	RuntimeError        Verdict = "RUNTIME_ERROR"
	TimeLimitExceeded   Verdict = "TIME_LIMIT_EXCEEDED"
	MemoryLimitExceeded Verdict = "MEMORY_LIMIT_EXCEEDED"
	WrongAnswer         Verdict = "WRONG_ANSWER"
	JudgementFailed     Verdict = "JUDGEMENT_FAILED"
	PresentationError   Verdict = "PRESENTATION_ERROR"
	Fail                Verdict = "FAIL"
	Points              Verdict = "POINTS"
	PartialResult       Verdict = "PARTIAL_RESULT"
	UnexpectedEof       Verdict = "UNEXPECTED_EOF"
	Dirt                Verdict = "DIRT"
)
