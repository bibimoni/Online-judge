package store

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg"
)

type StoreService interface {
	Get(id string) (pkg.Language, error)
	Register(l pkg.Language)
	List() []pkg.Language
	Contains(id string) bool
}

var DefaultStore StoreService = nil
