package contestservice

import (
	"context"
)

type ContestService interface {
	Create(ctx context.Context, author string) (string, error)
	AddPeople(ctx context.Context, contestId string, peopleType string, username string) error
	RemovePeople(ctx context.Context, contestId string, peopleType string, username string) error
}
