package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/benmanns/goworker"
)

var (
	// get flag for woker or publisher
	isRunWoker = flag.Bool("worker", false, "run worker")
	pushType   = flag.String("t", "", "enqueue type")
	sleepTime  = flag.Int("sleep", 0, "sleep job's sleep time")
)

func init() {
	settings := goworker.WorkerSettings{
		URI:            "redis://localhost:6379/",
		Connections:    100,
		Queues:         []string{"email", "sleep", "error", "panic", "schedule"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "sample-ns:",
		Interval:       5.0,
	}
	goworker.SetSettings(settings)
	goworker.Register("MyClass", myFunc)
}

func main() {
	flag.Parse()
	if *isRunWoker {
		runWoker()
	} else {
		push()
	}
}

func myFunc(queue string, args ...interface{}) error {
	log.Printf("From %s, %v\n", queue, args)
	var err error
	switch queue {
	case "email":
		err = doEmail(args)
	case "error":
		err = doError(args)
	case "panic":
		doPanic(args)
	case "schedule":
	default:
		panic(fmt.Sprintf("unknown queue: %s", queue))
	}
	return err
}

func runWoker() {
	log.Println("start job worker")
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}

func push() {
	switch *pushType {
	case "email":
		err := goworker.Enqueue(&goworker.Job{
			Queue: "email",
			Payload: goworker.Payload{
				Class: "MyClass",
				Args:  []interface{}{"sample@example.com", "test mail subject", "this is sample email"},
			},
		})
		if err != nil {
			panic(err)
		}
	case "error":
		err := goworker.Enqueue(&goworker.Job{
			Queue: "error",
			Payload: goworker.Payload{
				Class: "MyClass",
				Args:  []interface{}{},
			},
		})
		if err != nil {
			panic(err)
		}
	case "panic":
		err := goworker.Enqueue(&goworker.Job{
			Queue: "panic",
			Payload: goworker.Payload{
				Class: "MyClass",
				Args:  []interface{}{},
			},
		})
		if err != nil {
			panic(err)
		}
	case "schedule":
	default:
		log.Println("unsupportted push type:", *pushType)
	}
}

func doEmail(args []interface{}) error {
	log.Println("sending email...", args)
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		log.Println("sending...")
	}
	log.Println("done!!")
	return nil
}

func doError(args []interface{}) error {
	log.Println("starting error process.... ", args)
	return fmt.Errorf("error has occurred")
}

func doPanic(args []interface{}) {
	panic("manual panic")
}
