# workers

Got a lot of jobs to be done as fast as possible? This lib has got your back! It helps you run jobs concurrently with ease.

It creates goroutines to run any task you'd like and while there is nothing to do it waits.

## Options

You can make 2 configurations:

- Total Workers: the total goroutines that will start and wait to do the job. If not defined it defaults to 10.
- Buffer Channel: configures the channel that sends the jobs to the workers. If not defined it defaults to an unbuffered channel.

## How to use it

Call the `Start` function passing a context with cancel or timeout, that way the workers will stop when the context is canceled or if it times out. You can also pass options to configure the total number of workers and the size of the buffer channel.

WARNING: failing to pass a context that cancels makes the goroutines run forever until the end of the process.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
opt := workers.Options{
    TotalWorkers: 20, // it will spawn 20 goroutines
    BufferChannel: 0, // the channel will be unbuffered
}
workers.Start(ctx, opt)
```

You can also just call `Start` without options, in this case, it will start 10 goroutines and the channel will be unbuffered.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
workers.Start(ctx) //10 goroutines and unbuffered channel
```

After starting the workers you can call `DoJob` for the workers to do anything you need. This function receives the interface `Worker` so any struct that implements `Work` can be used here.

```go
type WorkJob struct{
    a int
    b int
}

func (w WorkJob) Work() {
    c := w.a + w.b
    log.Println("worker calculated a+b", c)
}

w := WorkJob {
    a: 2,
    b: 3,
}
err := workers.DoJob(w)
```

Another way is to use the `Job` type as follows. This type implements the interface `Worker`.

```go
var a, b = 2, 3
var j workers.Job = func() {
    c := a+b
    log.Println("worker calculated a+b", c)
}
err := workers.DoJob(j)
```

Just like magic, your job will be executed by one of the available workers.
