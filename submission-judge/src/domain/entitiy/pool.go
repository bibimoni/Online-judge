package domain

type Pool struct {
	Isolates chan *Isolate
}
