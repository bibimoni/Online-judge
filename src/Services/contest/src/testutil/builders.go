package testutil

import (
	"time"

	"contest/src/domain/entity"
	"github.com/google/uuid"
)

// ContestBuilder provides a fluent interface for creating test contests
type ContestBuilder struct {
	contest entity.Contest
}

// NewContestBuilder creates a new contest builder with sensible defaults
func NewContestBuilder() *ContestBuilder {
	now := time.Now()
	return &ContestBuilder{
		contest: entity.Contest{
			ID:          uuid.New().String(),
			Name:        "Test Contest",
			Description: "Test Description",
			StartTime:   now.Add(24 * time.Hour),
			EndTime:     now.Add(29 * time.Hour),
			ScoringType: entity.ScoringTypeICPC,
			Status:      entity.StatusDraft,
			Rated:       false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

func (b *ContestBuilder) WithID(id string) *ContestBuilder {
	b.contest.ID = id
	return b
}

func (b *ContestBuilder) WithName(name string) *ContestBuilder {
	b.contest.Name = name
	return b
}

func (b *ContestBuilder) WithDescription(desc string) *ContestBuilder {
	b.contest.Description = desc
	return b
}

func (b *ContestBuilder) WithTimeRange(start, end time.Time) *ContestBuilder {
	b.contest.StartTime = start
	b.contest.EndTime = end
	return b
}

func (b *ContestBuilder) WithScoringType(scoringType entity.ScoringType) *ContestBuilder {
	b.contest.ScoringType = scoringType
	return b
}

func (b *ContestBuilder) WithStatus(status entity.ContestStatus) *ContestBuilder {
	b.contest.Status = status
	return b
}

func (b *ContestBuilder) WithRated(rated bool) *ContestBuilder {
	b.contest.Rated = rated
	return b
}

func (b *ContestBuilder) WithRatingPolicy(policy entity.RatingPolicy) *ContestBuilder {
	b.contest.RatingPolicy = policy
	return b
}

func (b *ContestBuilder) AsICPC() *ContestBuilder {
	b.contest.ScoringType = entity.ScoringTypeICPC
	return b
}

func (b *ContestBuilder) AsIOI() *ContestBuilder {
	b.contest.ScoringType = entity.ScoringTypeIOI
	return b
}

func (b *ContestBuilder) AsOngoing() *ContestBuilder {
	now := time.Now()
	b.contest.StartTime = now.Add(-1 * time.Hour)
	b.contest.EndTime = now.Add(4 * time.Hour)
	b.contest.Status = entity.StatusOngoing
	return b
}

func (b *ContestBuilder) AsEnded() *ContestBuilder {
	now := time.Now()
	b.contest.StartTime = now.Add(-5 * time.Hour)
	b.contest.EndTime = now.Add(-1 * time.Hour)
	b.contest.Status = entity.StatusEnded
	return b
}

func (b *ContestBuilder) Build() entity.Contest {
	return b.contest
}

func (b *ContestBuilder) BuildPtr() *entity.Contest {
	return &b.contest
}

// ContestSubmissionBuilder provides a fluent interface for creating test submissions
type ContestSubmissionBuilder struct {
	submission entity.ContestSubmission
}

func NewContestSubmissionBuilder() *ContestSubmissionBuilder {
	now := time.Now()
	return &ContestSubmissionBuilder{
		submission: entity.ContestSubmission{
			ID:           uuid.New().String(),
			SubmissionID: uuid.New().String(),
			ContestID:    uuid.New().String(),
			Username:     "testuser",
			ProblemID:    "problem-A",
			SubmittedAt:  now,
			Verdict:      entity.VerdictPending,
			Points:       0,
			EvalStatus:   entity.EvalStatusPending,
			Ignored:      false,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}

func (b *ContestSubmissionBuilder) WithID(id string) *ContestSubmissionBuilder {
	b.submission.ID = id
	return b
}

func (b *ContestSubmissionBuilder) WithSubmissionID(submissionID string) *ContestSubmissionBuilder {
	b.submission.SubmissionID = submissionID
	return b
}

func (b *ContestSubmissionBuilder) WithContestID(contestID string) *ContestSubmissionBuilder {
	b.submission.ContestID = contestID
	return b
}

func (b *ContestSubmissionBuilder) WithUsername(username string) *ContestSubmissionBuilder {
	b.submission.Username = username
	return b
}

func (b *ContestSubmissionBuilder) WithProblemID(problemID string) *ContestSubmissionBuilder {
	b.submission.ProblemID = problemID
	return b
}

func (b *ContestSubmissionBuilder) WithVerdict(verdict entity.Verdict) *ContestSubmissionBuilder {
	b.submission.Verdict = verdict
	return b
}

func (b *ContestSubmissionBuilder) WithPoints(points float64) *ContestSubmissionBuilder {
	b.submission.Points = points
	return b
}

func (b *ContestSubmissionBuilder) AsAccepted() *ContestSubmissionBuilder {
	b.submission.Verdict = entity.VerdictAccepted
	b.submission.EvalStatus = entity.EvalStatusFinished
	b.submission.JudgedAt = time.Now()
	return b
}

func (b *ContestSubmissionBuilder) AsWrongAnswer() *ContestSubmissionBuilder {
	b.submission.Verdict = entity.VerdictWrongAnswer
	b.submission.EvalStatus = entity.EvalStatusFinished
	b.submission.JudgedAt = time.Now()
	return b
}

func (b *ContestSubmissionBuilder) AsIgnored() *ContestSubmissionBuilder {
	b.submission.Ignored = true
	return b
}

func (b *ContestSubmissionBuilder) Build() entity.ContestSubmission {
	return b.submission
}

func (b *ContestSubmissionBuilder) BuildPtr() *entity.ContestSubmission {
	return &b.submission
}

// ScoreboardBuilder provides a fluent interface for creating test scoreboards
type ScoreboardBuilder struct {
	snapshot entity.ScoreboardSnapshot
}

func NewScoreboardBuilder() *ScoreboardBuilder {
	return &ScoreboardBuilder{
		snapshot: entity.ScoreboardSnapshot{
			ContestID: uuid.New().String(),
			Kind:      entity.SnapshotKindLive,
			Rows:      []entity.ScoreboardRow{},
			CreatedAt: time.Now(),
		},
	}
}

func (b *ScoreboardBuilder) WithContestID(contestID string) *ScoreboardBuilder {
	b.snapshot.ContestID = contestID
	return b
}

func (b *ScoreboardBuilder) AsLive() *ScoreboardBuilder {
	b.snapshot.Kind = entity.SnapshotKindLive
	return b
}

func (b *ScoreboardBuilder) AsFinal() *ScoreboardBuilder {
	b.snapshot.Kind = entity.SnapshotKindFinal
	return b
}

func (b *ScoreboardBuilder) WithRow(row entity.ScoreboardRow) *ScoreboardBuilder {
	b.snapshot.Rows = append(b.snapshot.Rows, row)
	return b
}

func (b *ScoreboardBuilder) Build() entity.ScoreboardSnapshot {
	return b.snapshot
}

func (b *ScoreboardBuilder) BuildPtr() *entity.ScoreboardSnapshot {
	return &b.snapshot
}

// RejudgeJobBuilder provides a fluent interface for creating test rejudge jobs
type RejudgeJobBuilder struct {
	job entity.RejudgeJob
}

func NewRejudgeJobBuilder() *RejudgeJobBuilder {
	return &RejudgeJobBuilder{
		job: entity.RejudgeJob{
			ID:            uuid.New().String(),
			ContestID:     uuid.New().String(),
			SubmissionIDs: []string{},
			RequestedBy:   "admin",
			Status:        entity.RejudgeStatusPending,
			CreatedAt:     time.Now(),
		},
	}
}

func (b *RejudgeJobBuilder) WithContestID(contestID string) *RejudgeJobBuilder {
	b.job.ContestID = contestID
	return b
}

func (b *RejudgeJobBuilder) WithSubmissionIDs(ids ...string) *RejudgeJobBuilder {
	b.job.SubmissionIDs = ids
	return b
}

func (b *RejudgeJobBuilder) WithProblemID(problemID string) *RejudgeJobBuilder {
	b.job.ProblemID = &problemID
	return b
}

func (b *RejudgeJobBuilder) WithStatus(status entity.RejudgeStatus) *RejudgeJobBuilder {
	b.job.Status = status
	return b
}

func (b *RejudgeJobBuilder) Build() entity.RejudgeJob {
	return b.job
}

func (b *RejudgeJobBuilder) BuildPtr() *entity.RejudgeJob {
	return &b.job
}
