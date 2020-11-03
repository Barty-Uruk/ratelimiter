package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type commandStruct struct {
	Command string
	Arg     string
}

func main() {
	// wordPtr := flag.Int("rate", 1, "an int")
	// inflight := flag.Int("inflight", 1, "an int")
	// flag.Parse()
	// var param string
	// fmt.Scan(&param)
	command := os.Args[5]
	if command == "" {
		log.Fatal("command is empty")
	}
	paramChan := make(chan commandStruct)
	go scanCommandLineParams(paramChan, command)
	// for {
	// 	<-time.After(time.Second)
	// 	fmt.Println("1111")
	// }

	go execWithTimeout(1, 1, paramChan)
	// for i := 0; i < 50; i++ {
	// 	fmt.Println("======", command, i)
	// 	var param string
	// 	fmt.Scan(&param)
	// 	paramChan <- commandStruct{
	// 		Command: command,
	// 		Arg:     param,
	// 	}

	// }
	// fmt.Println(*wordPtr)
	// fmt.Println(*inflight)
	time.Sleep(5 * time.Second)
}

func execWithTimeout(timeOut time.Duration, limit int, commands chan commandStruct) {
	for i := 0; i <= limit; i++ {
		go func() {
			for comm := range commands {
				cmd := exec.Command(comm.Command, comm.Arg)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					log.Fatalf("cmd.Run() failed with %s\n", err)
				}
				time.Sleep(timeOut)
			}
		}()
	}

}

func scanCommandLineParams(commChan chan commandStruct, inputCommand string) {
	for i := 0; i < 150; i++ {
		fmt.Println("======", inputCommand)
		var param string
		fmt.Scan(&param)
		if param == "" {
			break
		}
		commChan <- commandStruct{
			Command: inputCommand,
			Arg:     param,
		}

	}
}
