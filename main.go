package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/panachainy/poc-golang-ants/background_job"
)

var wg sync.WaitGroup

func main() {
	bgj, err := background_job.NewAntsBackgroundJob(10)
	if err != nil {
		panic(err)
	}

	runTimes := 100
	wg.Add(runTimes)

	for i := 0; i < runTimes; i++ {
		taskId := i + 1
		err := bgj.Submit(
			withDurationLog(func() {
				defer wg.Done()
				exampleTask()
			},
				strconv.Itoa(taskId),
			),
		)
		if err != nil {
			fmt.Printf("Error submitting task %d: %v\n", taskId, err)
		}
		fmt.Printf("test after submit %v\n", taskId)
	}

	wg.Wait()

	// TODO: get return from task to do next step

	fmt.Println("done")
}

func exampleTask() {
	s := time.Second
	time.Sleep(s)
}

func withDurationLog(fn func(), msg string) func() {
	return func() {
		start := time.Now()
		fn()
		duration := time.Since(start)
		fmt.Printf("%s completed in: %v\n", msg, duration)
	}
}
