# go-worker-example goworker
example for https://github.com/benmanns/goworker

# Worker
```
$ go run main.go -worker true
```
```
2019/06/12 20:13:47 start job worker
```

# Enqueuer
## Email sample
```
$ go run main.go -t email
```
```
2019/06/12 20:16:04 From email, [sample@example.com test mail subject this is sample email]
2019/06/12 20:16:04 sending email... [sample@example.com test mail subject this is sample email]
2019/06/12 20:16:05 sending...
2019/06/12 20:16:06 sending...
2019/06/12 20:16:07 sending...
2019/06/12 20:16:08 sending...
2019/06/12 20:16:09 sending...
2019/06/12 20:16:09 done!!
```

## error
```
$ go run main.go -t error
```
```
2019/06/12 20:16:47 From error, []
2019/06/12 20:16:47 starting error process....  []
```
Nothing happens...

## panic
```
$ go run main.go -t panic
```
```
2019/06/12 20:17:34 From panic, []
```
Nothing happens...
