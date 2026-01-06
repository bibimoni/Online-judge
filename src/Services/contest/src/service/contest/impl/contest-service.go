package contestserviceimpl

import (
	"contest/src/domain/repository/contest"
	service "contest/src/service/contest"
	"context"
)

type ContestServiceImpl struct {
	contestRepo contestrepo.ContestRepository
}

func NewContestService(contestRepo contestrepo.ContestRepository) service.ContestService {
	return &ContestServiceImpl{
		contestRepo: contestRepo,
	}
}

func (s *ContestServiceImpl) Create(ctx context.Context, author, contestname string) (string, error) {
	return s.contestRepo.Create(ctx, author, contestname)
}

func (s *ContestServiceImpl) AddPeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error {
	return s.contestRepo.AddPeople(ctx, contestId, peopleType, username)
}

func (s *ContestServiceImpl) RemovePeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error {
	return s.contestRepo.RemovePeople(ctx, contestId, peopleType, username)
}
