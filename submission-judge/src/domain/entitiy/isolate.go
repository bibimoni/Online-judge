package domain

import (
	"github.com/rs/zerolog"
)

// This will be a sandbox implementation, it will call isolate command.
// It needs isolate installed in the system to work
type Isolate struct {
	ID       int
	BoxDir   string
	MetaFile string
	Logger   *zerolog.Logger
	Inited   bool
}
