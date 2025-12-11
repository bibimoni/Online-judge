package contest

import (
	"context"
)

type ContestService interface {
	Create(ctx context.Context, authorId uint64) (string, error)
	AddPeople(contestId string, peopleType string, userId uint64) error
	RemovePeople(contestId string, peopleType string, userId uint64) error
}
