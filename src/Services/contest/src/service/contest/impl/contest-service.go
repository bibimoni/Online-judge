package contestserviceimpl

import (
	repository "contest/src/domain/repository/contest"
	service "contest/src/service/contest"
	"context"
)

type ContestServiceImpl struct {
	contestRepo repository.ContestRepository
}

func NewContestService(contestRepo repository.ContestRepository) service.ContestService {
	return &ContestServiceImpl{
		contestRepo: contestRepo,
	}
}

func (s *ContestServiceImpl) Create(ctx context.Context, author string) (string, error) {
	return s.contestRepo.Create(ctx, author)
}

func (s *ContestServiceImpl) AddPeople(ctx context.Context, contestId string, peopleType string, username string) error {
	return s.contestRepo.AddPeople(ctx, contestId, peopleType, username)
}

func (s *ContestServiceImpl) RemovePeople(ctx context.Context, contestId string, peopleType string, username string) error {
	return s.contestRepo.RemovePeople(ctx, contestId, peopleType, username)
}
