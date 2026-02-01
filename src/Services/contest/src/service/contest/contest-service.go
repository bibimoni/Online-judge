package contestservice

import (
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	"context"
	"errors"
)

type ContestService interface {
	Create(ctx context.Context, author, contestname string) (string, error)
	AddPeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string, participantType domain.ParticipantType) error
	RemovePeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error
	ChangeProblems(ctx context.Context, contestId string, problemIds []string, shortName []string) error
	IsProblemInActiveCOntest(ctx context.Context, problemId uint64) (bool, *domain.Contest, error)
}

var ErrProblemNotFound = errors.New("problem not found")
