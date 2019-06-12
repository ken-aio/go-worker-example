# go-worker-example work
example for https://github.com/gocraft/work

# Worker
```
$ go run main.go -worker true
```
```
2019/06/12 17:10:07 begin worker startup...
2019/06/12 17:10:07 starting worker process now
```

# Enqueuer
## Email sample
```
$ go run main.go -t email
```
```
2019/06/12 17:11:03 starting job:  email
2019/06/12 17:11:03 sending email... map[address:sample@example.com body:this is sample email subject:test mail subject]
2019/06/12 17:11:04 sending...
2019/06/12 17:11:05 sending...
2019/06/12 17:11:06 sending...
2019/06/12 17:11:07 sending...
2019/06/12 17:11:08 sending...
2019/06/12 17:11:08 done!!
```

## error and retry sample
```
$ go run main.go -t error
```
```
2019/06/12 17:12:23 starting job:  error
2019/06/12 17:12:23 starting error process.... number of errors is  0
2019/06/12 17:13:23 starting job:  error
2019/06/12 17:13:23 starting error process.... number of errors is  1
2019/06/12 17:14:48 starting job:  error
2019/06/12 17:14:48 starting error process.... number of errors is  2
```

## panic and retry sample
```
$ go run main.go -t panic
```
```
2019/06/12 17:15:13 starting panic process.... number of errors is  0
ERROR: runJob.panic - manual panic
2019/06/12 17:16:13 starting job:  panic
2019/06/12 17:16:13 starting panic process.... number of errors is  1
ERROR: runJob.panic - manual panic
2019/06/12 17:16:59 starting job:  panic
2019/06/12 17:16:59 starting panic process.... number of errors is  2
ERROR: runJob.panic - manual panic
2019/06/12 17:17:53 starting job:  error
2019/06/12 17:17:53 starting error process.... number of errors is  3
```

## schedule sample
Run job after -sleep time sec  
```
$ go run main.go -t schedule -sleep 10
2019/06/12 17:20:45 exec push. type: schedule
2019/06/12 17:20:45 reserved schedule job after 10 sec
```
```
2019/06/12 17:20:55 starting job:  schedule
2019/06/12 17:20:55 done!!
```
