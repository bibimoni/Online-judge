package impl

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

func (s *ContestServiceImpl) Create(ctx context.Context, authorId uint64) (string, error) {
	return s.contestRepo.Create(ctx, authorId)
}

func (s *ContestServiceImpl) AddPeople(contestId string, peopleType string, userId uint64) error {
	return s.contestRepo.AddPeople(contestId, peopleType, userId)
}

func (s *ContestServiceImpl) RemovePeople(contestId string, peopleType string, userId uint64) error {
	return s.contestRepo.RemovePeople(contestId, peopleType, userId)
}
