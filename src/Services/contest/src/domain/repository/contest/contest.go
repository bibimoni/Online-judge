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

	AddPeople(ctx context.Context, contestId string, peopleType PeopleType, username string) error
	RemovePeople(ctx context.Context, contestId string, peopleType PeopleType, username string) error
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
