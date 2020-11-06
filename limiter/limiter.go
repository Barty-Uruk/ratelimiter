package limiter

import (
	"sync"
)

type Limiter struct {
	Inflight  int
	Executors []*Executor
}

func NewLimiter(rate, inflight int, command string, args []string) *Limiter {
	limiter := Limiter{Inflight: inflight}
	for i := 0; i < inflight; i++ {
		limiter.Executors = append(limiter.Executors, NewExecutor(rate, command, args))
	}
	return &limiter
}

func (l *Limiter) Exec(input []string) {
	for _, exe := range l.Executors {
		go exe.Start()
	}

	var wg sync.WaitGroup
	for i := 0; i < len(input); i++ {
		exe := l.Executors[i%len(l.Executors)]
		wg.Add(1)
		exe.Exec(Input{
			Data: input[i],
			Wg:   &wg,
		})
	}

	for _, exe := range l.Executors {
		exe.Stop()
	}

	wg.Wait()
}
