package impl

import (
	"errors"
	"sync"

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
	currentCount   int
	mu             sync.Mutex
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

	// // Init all isolate
	// for i := cfg.Judge.IDOffset; i < (cfg.Judge.IDOffset + cfg.Judge.Amount); i++ {
	// 	newIsolate, err := is.NewIsolate(i)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	err = is.Init(newIsolate)
	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }
	// 	newPool.Put(newIsolate)
	// }

	config.GetLogger().Info().Msgf("Finished initialized pool service")

	return newPool, nil
}

func (ps *PoolServiceImpl) Get() (*domain.Isolate, error) {
	select {
	case i, ok := <-ps.pool.Isolates:
		if !ok {
			return nil, errors.New("Channel is closed")
		}
		config.GetLogger().Info().Msgf("Isolate available, use it")
		return i, nil
	default:
		cfg, err := config.Load()
		config.GetLogger().Info().Msgf("Offset: %d, amount: %d", cfg.Judge.IDOffset, cfg.Judge.Amount)
		if err != nil {
			return nil, err
		}
		ps.mu.Lock()
		if ps.currentCount < cfg.Judge.Amount {
			id := cfg.Judge.IDOffset + ps.currentCount
			newIsolate, err := ps.isolateService.NewIsolate(id)
			if err != nil {
				return nil, err
			}
			err = ps.isolateService.Init(newIsolate)
			ps.currentCount += 1
			ps.mu.Unlock()
			return newIsolate, nil
		}
		ps.mu.Unlock()
		i, ok := <-ps.pool.Isolates
		if !ok {
			return nil, errors.New("Channel is closed")
		}
		return i, nil

	}
}

func (ps *PoolServiceImpl) Put(i *domain.Isolate) {
	ps.pool.Isolates <- i
}

func (ps *PoolServiceImpl) Len() int {
	return len(ps.pool.Isolates)
}
