package poolservice

import (
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
)

type PoolService interface {
	Get() (*domain.Isolate, error)
	Put(i *domain.Isolate)
	Len() int
}
