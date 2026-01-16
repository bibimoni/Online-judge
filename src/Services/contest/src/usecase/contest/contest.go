package contestusecase

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	"context"
	"time"
)

type ContestOutput struct {
	Status string `json:"status"`
}

type ContestInteractor interface {
	EditContest(ctx context.Context, input *EditContestInput) (*EditContestOutput, error)
	CreateContest(ctx context.Context, input *CreateContestInput) (*CreateContestOutput, error)
	PatchContest(ctx context.Context, input *PatchContestInput) (*PatchContestOutput, error)
	ManageContestProblems(ctx context.Context, input *ManageContestProblemsInput) (*ManageContestProblemsOutput, error)
	GetContestById(ctx context.Context, input *GetContestByIdInput) (*GetContestByIdOutput, error)
	GetAllContests(ctx context.Context, input *GetContestsInput) (*GetContestsOutput, error)
}

type (
	GetContestsOutput struct {
		Contests []*domain.Contest `json:"contests"`
	}
	GetContestByIdOutput struct {
		Contest *domain.Contest `json:"contest"`
	}
	GetContestByIdInput struct {
		ContestId string
		GetContestsInput
	}
	GetContestsInput struct {
		Username      string
		Role          string
		Authenticated bool
	}
	EditContestInput struct {
		EditType   EditType               `json:"_"` // From Query Param
		ContestId  string                 `json:"contest_id" validate:"required"`
		PeopleType contestrepo.PeopleType `json:"people_type" validate:"required"`
		Target     string                 `json:"target" validate:"required"`
		Username   string                 `json:"username" validate:"required"` // Auto injected from gateway
		UserRole   string                 `json:"user_role,omitempty"`
	}

	EditContestOutput struct {
		common.StatusOK
	}

	CreateContestInput struct {
		UserRole string // In header X-User-Role
		Username string // Creator of the contest, auto injected from gateway
		Name     string
	}

	CreateContestOutput struct {
		ContestId string `json:"contest_id"`
	}

	PatchContestOutput struct {
		common.StatusOK
	}

	PatchContestInput struct {
		ContestId            string                       `json:"contest_id" validate:"required"`
		Username             string                       `json:"username,omitempty" validate:"required"`
		UserRole             string                       `json:"user_role,omitempty"`
		Description          *string                      `validate:"max=500" json:"description,omitempty"`
		ScoreboardVisibility *domain.ScoreboardVisibility `json:"scoreboard_visibility,omitempty"`
		StartTime            *time.Time                   `json:"start_time,omitempty"` // RFC 3339 format
		EndTime              *time.Time                   `json:"end_time,omitempty"`   // RFC 3339 format
		ContestRule          *PatchContestRuleInput       `json:"contest_rule,omitempty"`
		FinalizeAt           *time.Time                   `json:"finalize_at,omitempty"`
		RejudgeWindowEnd     *time.Time                   `json:"rejudge_window_end,omitempty"`
	}

	PatchContestRuleInput struct {
		ScoringType                     *domain.ScoringType `json:"scoring_type,omitempty"`
		PenaltyMinutes                  *uint16             `json:"penalty_minutes,omitempty"`
		FreezeStartTime                 *time.Time          `json:"freeze_start_time"`
		FreezeTime                      *time.Duration      `json:"freeze_time,omitempty"`
		MaxAllowedSubmissionsPerProblem *int16              `json:"max_allowed_submissions_per_problem,omitempty"`
	}

	ManageContestProblemsInput struct {
		ContestId  string   `json:"contest_id,omitempty"`
		Username   string   `json:"username,omitempty"`
		ProblemIds []string `json:"problem_ids,omitempty"`
		ShortNames []string `json:"short_names,omitempty"`
	}

	ManageContestProblemsOutput struct {
		common.StatusOK
	}
)

type EditType string

const (
	AddPeople    EditType = "add_people"
	RemovePeople EditType = "remove_people"
)
