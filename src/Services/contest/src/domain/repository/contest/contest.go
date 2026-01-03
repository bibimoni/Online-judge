package contestrepo

import (
	domain "contest/src/domain/entity"
	"context"
	"errors"
)

type ContestRepository interface {
	GetById(ctx context.Context, contestId string) (*domain.Contest, error)
	Create(ctx context.Context, creator string) (string, error)

	// AddContestant(contestId string, userId uint64) error

	AddPeople(ctx context.Context, contestId string, peopleType string, username string) error
	RemovePeople(ctx context.Context, contestId string, peopleType string, username string) error
}

var NoContestFound = errors.New("no contest found")

const (
	Author     string = "AUTHOR"
	Admin      string = "ADMIN"
	Tester     string = "TESTER"
	Contestant string = "CONTESTANT"
)

var ContestPeopple = []string{Author, Admin, Tester, Contestant}
