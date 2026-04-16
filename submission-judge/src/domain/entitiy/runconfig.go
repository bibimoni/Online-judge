package domain

import (
	"io"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
)

// This will be the arguments passed to isolate
type RunConfig struct {
	TimeLimit        time.Duration
	MemoryLimit      memory.Memory
	WorkingDirectory string
	DirectoryMaps    []DirectoryMap
	Env              []string
	InheritEnv       bool
	MaxProcesses     int
	Input            string
	Output           string
	Meta             bool
	Args             []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type DirectoryMap struct {
	Inside  string
	Outside string
	Options []DirectoryMapOption
}

type DirectoryMapOption string

const (
	AllowSpecial   DirectoryMapOption = "dev"
	MountFS        DirectoryMapOption = "fs"
	Maybe          DirectoryMapOption = "maybe"
	NoExec         DirectoryMapOption = "noexec"
	AllowReadWrite DirectoryMapOption = "rw"
	Temporary      DirectoryMapOption = "tmp"
)
