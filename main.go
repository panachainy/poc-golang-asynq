package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/panachainy/poc-golang-asynq/tasks"
)

var wg sync.WaitGroup

func main() {
	const redisAddr = "127.0.0.1:6379"

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	info, err = client.Enqueue(task, asynq.ProcessIn(24*time.Hour))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	task, err = tasks.NewImageResizeTask("https://example.com/myassets/image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// bgj, err := background_job.NewAntsBackgroundJob(10)
	// if err != nil {
	// 	panic(err)
	// }

	// runTimes := 100
	// wg.Add(runTimes)

	// for i := 0; i < runTimes; i++ {
	// 	taskId := i + 1
	// 	err := bgj.Submit(
	// 		withDurationLog(func() {
	// 			defer wg.Done()
	// 			exampleTask()
	// 		},
	// 			strconv.Itoa(taskId),
	// 		),
	// 	)
	// 	if err != nil {
	// 		fmt.Printf("Error submitting task %d: %v\n", taskId, err)
	// 	}
	// 	fmt.Printf("test after submit %v\n", taskId)
	// }

	// wg.Wait()

	// // TODO: get return from task to do next step

	// fmt.Println("done")
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
