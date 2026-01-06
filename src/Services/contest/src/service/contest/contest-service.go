package contestservice

import (
	contestrepo "contest/src/domain/repository/contest"
	"context"
)

type ContestService interface {
	Create(ctx context.Context, author, contestname string) (string, error)
	AddPeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error
	RemovePeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error
}
