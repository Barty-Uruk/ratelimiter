package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/Barty-Uruk/ratelimiter/limiter"
)

const (
	defaultRate     = 1
	defaultInflight = 1
)

func main() {
	rate := flag.Int("rate", defaultRate, "command rate per a second")
	inflight := flag.Int("inflight", defaultInflight, "maximum parallel commands allowed")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("you must provide a command")
	}
	command, args := args[0], args[1:]

	var commands []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comm := scanner.Text()
		if comm == "" {
			break
		}
		commands = append(commands, comm)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	lim := limiter.NewLimiter(*rate, *inflight, command, args)
	lim.Exec(commands)
}
