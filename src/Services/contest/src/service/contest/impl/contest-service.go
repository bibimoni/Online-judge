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

func (s *ContestServiceImpl) AddPeople(contestId string, peopleType string, username string) error {
	return s.contestRepo.AddPeople(contestId, peopleType, username)
}

func (s *ContestServiceImpl) RemovePeople(contestId string, peopleType string, username string) error {
	return s.contestRepo.RemovePeople(contestId, peopleType, username)
}
