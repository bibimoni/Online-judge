package contestrepo

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

type ContestRepository interface {
	GetById(ctx context.Context, contestId string) (*domain.Contest, error)
	Create(ctx context.Context, creator string, contestname string) (string, error)
	ReplaceOne(ctx context.Context, contestId string, updatedContest *domain.Contest) error
	UpdateOne(ctx context.Context, contestId string, updateData map[string]any) error
	CanCreateContest(ctx context.Context, role string) bool

	// AddContestant(contestId string, userId uint64) error

	AddPeople(ctx context.Context, contestId string, peopleType PeopleType, username string, participantType domain.ParticipantType) error
	RemovePeople(ctx context.Context, contestId string, peopleType PeopleType, username string) error

	UpdateProblems(ctx context.Context, contestId string, problems []*domain.ContestProblem) error
	GetProblemByLabels(ctx context.Context, contestId string, problemLabels []string) ([]*domain.ContestProblem, error)

	ListAllContestsWithAuth(ctx context.Context, username, role string) ([]*domain.Contest, error)
	ListAllPublicContests(ctx context.Context) ([]*domain.Contest, error)

	GetContestsByProblemId(ctx context.Context, problemId uint64) ([]*domain.Contest, error)
}

var ErrNoContestFound = errors.New("no contest found")

type PeopleType string

const (
	Author     PeopleType = "AUTHOR"
	Admin      PeopleType = "ADMIN"
	Tester     PeopleType = "TESTER"
	Contestant PeopleType = "CONTESTANT"
)

var ContestPeopple = []PeopleType{Author, Admin, Tester, Contestant}
