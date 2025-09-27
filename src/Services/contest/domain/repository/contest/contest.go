package repository

import (
	domain "contest/domain/entity"
	"context"
)

type ContestRepository interface {
	GetById(contestId string) (domain.Contest, error)
	Create(ctx context.Context, author uint64) (string, error)

	// AddContestant(contestId string, userId uint64) error

	AddPeople(contestId string, peopleType string, userId uint64) error
	RemovePeople(contestId string, peopleType string, userId uint64) error
	EditProblemset(contestId string, newProblemset []uint64) error
}
