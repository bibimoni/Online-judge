package impl

import (
	"errors"

	"fmt"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/impl"
	poolservice "github.com/bibimoni/Online-judge/submission-judge/src/service/pool"
)

type PoolServiceImpl struct {
	pool           *domain.Pool
	isolateService isolateservice.IsolateService
}

func NewPoolSerivce() (poolservice.PoolService, error) {
	poolService, err := NewPoolServiceImpl()
	if err != nil {
		return nil, fmt.Errorf("Error when create new Pool %v", err)
	}
	return poolService, nil
}

func NewPoolServiceImpl() (*PoolServiceImpl, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	is, err := impl.NewIsolateService()
	if err != nil {
		return nil, err
	}

	newPool := &PoolServiceImpl{
		pool: &domain.Pool{
			Isolates: make(chan *domain.Isolate, cfg.Judge.Amount),
		},
		isolateService: is,
	}

	log := config.GetLogger()
	log.Info().Msgf("Offset: %d, amount: %d", cfg.Judge.IDOffset, cfg.Judge.Amount)
	// Init all isolate
	for i := cfg.Judge.IDOffset; i < (cfg.Judge.IDOffset + cfg.Judge.Amount); i++ {
		newIsolate, err := is.NewIsolate(i)
		if err != nil {
			return nil, err
		}
		err = is.Init(newIsolate)
		// if err != nil {
		// 	return nil, err
		// }
		newPool.Put(newIsolate)
	}

	log.Info().Msgf("Finished initialized pool service")

	return newPool, nil
}

func (ps *PoolServiceImpl) Get() (*domain.Isolate, error) {
	i, ok := <-ps.pool.Isolates
	if !ok {
		return nil, errors.New("Channel is closed")
	}
	return i, nil
}

func (ps *PoolServiceImpl) Put(i *domain.Isolate) {
	ps.pool.Isolates <- i
}

func (ps *PoolServiceImpl) Len() int {
	return len(ps.pool.Isolates)
}
