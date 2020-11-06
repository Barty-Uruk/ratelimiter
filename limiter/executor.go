package limiter

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	outMtx sync.Mutex
)

const (
	rateInterval    = time.Second
	noStdinArgIndex = -1
)

type Input struct {
	Data string
	Wg   *sync.WaitGroup
}

type Executor struct {
	Command string
	Args    []string

	input           chan Input
	limiter         *rate.Limiter
	replaceArgIndex int
}

func NewExecutor(rps int, command string, args []string) *Executor {
	exe := Executor{
		Command: command,
		Args:    args,

		input:   make(chan Input),
		limiter: rate.NewLimiter(rate.Every(rateInterval/time.Duration(rps)), 1),
	}
	exe.defineReplaceArgIndex()
	return &exe
}

func (e *Executor) Start() {
	for in := range e.input {
		e.execCommand(in)
	}
}

func (e *Executor) Exec(input Input) {
	e.input <- input
}

func (e *Executor) Stop() {
	close(e.input)
}

func (e *Executor) execCommand(input Input) {
	defer input.Wg.Done()

	e.wait()

	cmd := exec.Command(e.Command, e.replaceStdin(input.Data)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run a command: %v", err)
	}

	outMtx.Lock()
	fmt.Print(string(out))
	outMtx.Unlock()
}

func (e *Executor) wait() {
	now := time.Now()
	ticket := e.limiter.ReserveN(now, 1)
	if ticket.OK() {
		delay := ticket.DelayFrom(now)
		time.Sleep(delay)
	}
}

func (e *Executor) replaceStdin(input string) []string {
	if e.replaceArgIndex == noStdinArgIndex {
		return e.Args
	}

	args := make([]string, len(e.Args))
	copy(args, e.Args)
	args[e.replaceArgIndex] = input
	return args
}

func (e *Executor) defineReplaceArgIndex() {
	e.replaceArgIndex = noStdinArgIndex
	for i, arg := range e.Args {
		if arg == "{}" {
			e.replaceArgIndex = i
			return
		}
	}
}
