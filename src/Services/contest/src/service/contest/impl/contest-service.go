package contestserviceimpl

import (
	domain "contest/src/domain/entity"
	"contest/src/domain/repository/contest"
	"contest/src/infrastructure/config"
	contestservice "contest/src/service/contest"
	service "contest/src/service/contest"
	problemservice "contest/src/service/problem"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type ContestServiceImpl struct {
	contestRepo    contestrepo.ContestRepository
	problemService problemservice.ProblemService
}

func NewContestService(
	contestRepo contestrepo.ContestRepository,
	problemService problemservice.ProblemService,
) service.ContestService {
	return &ContestServiceImpl{
		contestRepo:    contestRepo,
		problemService: problemService,
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

func (s *ContestServiceImpl) ChangeProblems(ctx context.Context, contestId string, problemIds []string, shortNames []string) error {
	// chain of problems
	if len(problemIds) != len(shortNames) {
		return fmt.Errorf("length of problemIds and shortNames must be equal")
	}
	problems := make([]*domain.ContestProblem, len(problemIds))

	g, gCtx := errgroup.WithContext(ctx)
	for i := range len(problemIds) {
		g.Go(func() error {
			if gCtx.Err() != nil {
				return gCtx.Err()
			}

			config.GetLogger().Info().Msgf("Fetching problem %s from problem service", problemIds[i])
			problem, err := s.problemService.Get(gCtx, problemIds[i])
			if err != nil {
				config.GetLogger().Error().Err(err).Msgf("failed to get problem %s from problem service", problemIds[i])
				return contestservice.ErrProblemNotFound
			}

			problems[i] = &domain.ContestProblem{
				ProblemId: uint64(problem.ProblemId),
				Label:     shortNames[i],
				MaxPoints: s.problemService.GetProblemMaxScores(problem),
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	config.GetLogger().Info().Msgf("Updating problems for contest %s with %+v problems", contestId, problems)
	s.contestRepo.UpdateProblems(ctx, contestId, problems)
	return nil
}
