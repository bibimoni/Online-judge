package contestservice

import (
	"context"
)

type ContestService interface {
	Create(ctx context.Context, author string) (string, error)
	AddPeople(contestId string, peopleType string, username string) error
	RemovePeople(contestId string, peopleType string, username string) error
}
