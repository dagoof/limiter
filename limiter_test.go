package limiter

import (
	"fmt"
	"time"
)

func ExampleLimiter() {
	var (
		njobs   = 5
		limiter = New(2)
		jobs    = make(chan int)
		results = make(chan int)
	)

	work := func(l *Limiter, in <-chan int, out chan<- int) {
		n := <-in
		defer fmt.Println("done working", n)

		time.Sleep(time.Millisecond) // do some long-running work
		out <- n
		l.Done()
	}

	go func() {
		for i := 0; i < njobs; i++ {
			limiter.Wait()

			fmt.Println("spawning worker", i)
			go work(limiter, jobs, results)
			jobs <- i
		}
		limiter.WaitDone()
		close(results)
	}()

	for _ = range results {
	}

	// Output:
	// spawning worker 0
	// spawning worker 1
	// done working 0
	// spawning worker 2
	// done working 1
	// spawning worker 3
	// done working 2
	// spawning worker 4
	// done working 3
	// done working 4
}
