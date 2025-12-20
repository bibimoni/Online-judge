package testutil

import (
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/google/uuid"
)

// SubmissionBuilder provides fluent API for creating test submissions
type SubmissionBuilder struct {
	submission domain.Submission
}

func NewSubmissionBuilder() *SubmissionBuilder {
	return &SubmissionBuilder{
		submission: domain.Submission{
			ID:          uuid.New().String(),
			ProblemID:   "problem-123",
			UserID:      "user-123",
			Language:    "cpp",
			Code:        "int main() { return 0; }",
			Status:      domain.StatusPending,
			SubmittedAt: time.Now(),
		},
	}
}

func (b *SubmissionBuilder) WithID(id string) *SubmissionBuilder {
	b.submission.ID = id
	return b
}

func (b *SubmissionBuilder) WithProblemID(problemID string) *SubmissionBuilder {
	b.submission.ProblemID = problemID
	return b
}

func (b *SubmissionBuilder) WithUserID(userID string) *SubmissionBuilder {
	b.submission.UserID = userID
	return b
}

func (b *SubmissionBuilder) WithLanguage(lang string) *SubmissionBuilder {
	b.submission.Language = lang
	return b
}

func (b *SubmissionBuilder) WithCode(code string) *SubmissionBuilder {
	b.submission.Code = code
	return b
}

func (b *SubmissionBuilder) AsPending() *SubmissionBuilder {
	b.submission.Status = domain.StatusPending
	return b
}

func (b *SubmissionBuilder) AsJudging() *SubmissionBuilder {
	b.submission.Status = domain.StatusJudging
	return b
}

func (b *SubmissionBuilder) AsAccepted() *SubmissionBuilder {
	b.submission.Status = domain.StatusFinished
	b.submission.Verdict = domain.VerdictAccepted
	return b
}

func (b *SubmissionBuilder) AsWrongAnswer() *SubmissionBuilder {
	b.submission.Status = domain.StatusFinished
	b.submission.Verdict = domain.VerdictWrongAnswer
	return b
}

func (b *SubmissionBuilder) Build() domain.Submission {
	return b.submission
}

func (b *SubmissionBuilder) BuildPtr() *domain.Submission {
	return &b.submission
}

// JudgeResultBuilder for creating test judge results
type JudgeResultBuilder struct {
	result domain.JudgeResult
}

func NewJudgeResultBuilder() *JudgeResultBuilder {
	return &JudgeResultBuilder{
		result: domain.JudgeResult{
			SubmissionID:    uuid.New().String(),
			Verdict:         domain.VerdictAccepted,
			TestCasesPassed: 10,
			TotalTestCases:  10,
			ExecutionTime:   100,
			MemoryUsed:      1024,
		},
	}
}

func (b *JudgeResultBuilder) WithVerdict(verdict domain.Verdict) *JudgeResultBuilder {
	b.result.Verdict = verdict
	return b
}

func (b *JudgeResultBuilder) WithTestCases(passed, total int) *JudgeResultBuilder {
	b.result.TestCasesPassed = passed
	b.result.TotalTestCases = total
	return b
}

func (b *JudgeResultBuilder) Build() domain.JudgeResult {
	return b.result
}

// IsolateBuilder for creating test isolates
type IsolateBuilder struct {
	isolate domain.Isolate
}

func NewIsolateBuilder() *IsolateBuilder {
	return &IsolateBuilder{
		isolate: domain.Isolate{
			ID:      1,
			Inited:  true,
			Busy:    false,
		},
	}
}

func (b *IsolateBuilder) WithID(id int) *IsolateBuilder {
	b.isolate.ID = id
	return b
}

func (b *IsolateBuilder) AsBusy() *IsolateBuilder {
	b.isolate.Busy = true
	return b
}

func (b *IsolateBuilder) Build() domain.Isolate {
	return b.isolate
}
