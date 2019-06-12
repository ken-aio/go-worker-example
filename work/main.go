package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	// get flag for woker or publisher
	isRunWoker = flag.Bool("worker", false, "run worker")
	pushType   = flag.String("t", "", "enqueue type")
	sleepTime  = flag.Int("sleep", 0, "sleep job's sleep time")

	// Make a redis pool
	redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}

	// make an enqueuer with particular namespace
	enqueuer = work.NewEnqueuer("sample_namespace", redisPool)
)

const (
	queueEmail    = "email"
	queueError    = "error"
	queuePanic    = "panic"
	queueSleep    = "sleep"
	queueSchedule = "schedule"
)

// SampleContext context for sample
type SampleContext struct {
	MyName string
}

func main() {
	flag.Parse()
	if *isRunWoker {
		runWoker()
	} else {
		push()
	}
}

func runWoker() {
	log.Println("begin worker startup...")
	// Make a new pool. Arguments:
	// SampleContext{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "sample_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(SampleContext{}, 10, "sample_namespace", redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*SampleContext).Log)

	// make job
	pool.Job(queueEmail, (*SampleContext).SendMail)
	pool.Job(queueSleep, (*SampleContext).Sleep)
	pool.Job(queueSchedule, (*SampleContext).Schedule)

	// make job with options
	pool.JobWithOptions(queueError, work.JobOptions{Priority: 10, MaxFails: 5}, (*SampleContext).ErrorProcess)
	pool.JobWithOptions(queuePanic, work.JobOptions{Priority: 20, MaxFails: 3}, (*SampleContext).PanicProcess)

	// start job process
	pool.Start()

	log.Println("starting worker process now")
	// wait for signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	log.Println("quit worker process now... byebye")
	pool.Stop()
}

// Log is logging every job
func (c *SampleContext) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Println("starting job: ", job.Name)
	return next()
}

// SendMail is sending email for user
func (c *SampleContext) SendMail(job *work.Job) error {
	log.Println("sending email...", job.Args)
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		log.Println("sending...")
		job.Checkin(fmt.Sprintf("i = %d", i))
	}
	log.Println("done!!")
	return nil
}

// Sleep is sleep job
func (c *SampleContext) Sleep(job *work.Job) error {
	t := job.ArgInt64("time")
	log.Println("now sleeping", t, "sec...")
	time.Sleep(time.Duration(t) * time.Second)
	log.Println("done!!")
	return nil
}

// Schedule is scheduled job
func (c *SampleContext) Schedule(job *work.Job) error {
	log.Println("done!!")
	return nil
}

// ErrorProcess is too long sleep function
func (c *SampleContext) ErrorProcess(job *work.Job) error {
	log.Println("starting error process.... number of errors is ", job.Fails)
	return fmt.Errorf("error has occurred")
}

// PanicProcess is too long sleep function
func (c *SampleContext) PanicProcess(job *work.Job) error {
	log.Println("starting panic process.... number of errors is ", job.Fails)
	panic("manual panic")
}

func push() {
	log.Println("exec push. type:", *pushType)
	switch *pushType {
	case "email":
		_, err := enqueuer.Enqueue(queueEmail, work.Q{"address": "sample@example.com", "subject": "test mail subject", "body": "this is sample email"})
		if err != nil {
			panic(err)
		}
	case "sleep":
		_, err := enqueuer.Enqueue(queueSleep, work.Q{"time": *sleepTime})
		if err != nil {
			panic(err)
		}
	case "error":
		_, err := enqueuer.Enqueue(queueError, work.Q{})
		if err != nil {
			panic(err)
		}
	case "panic":
		_, err := enqueuer.Enqueue(queuePanic, work.Q{})
		if err != nil {
			panic(err)
		}
	case "schedule":
		log.Println("reserved schedule job after", *sleepTime, "sec")
		_, err := enqueuer.EnqueueIn(queueSchedule, int64(*sleepTime), work.Q{})
		if err != nil {
			panic(err)
		}
	default:
		log.Println("unsupportted push type:", *pushType)
	}
}
