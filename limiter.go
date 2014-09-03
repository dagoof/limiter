/*
Package limiter provides a singular type that allows you to control concurrent
execution by specifying a maximum number of tasks.

It allows you to then block until tasks complete before spawning new ones, and
also block on all tasks until they have completed.
*/
package limiter

import "sync"

// Limiter is a coordinator that waits on current tasks to complete before 
// allowing new ones to start.
type Limiter struct {
	cond        *sync.Cond
	wg          *sync.WaitGroup
	active, max int
}

// New creates a Limiter with a maximum number of concurrent tasks.
func New(max int) *Limiter {
	return &Limiter{
		sync.NewCond(&sync.Mutex{}),
		&sync.WaitGroup{},
		0,
		max,
	}
}

func (l *Limiter) ready() bool {
	return l.active < l.max
}

func (l *Limiter) done() {
	l.active--
	l.wg.Done()
}

func (l *Limiter) add() {
	l.active++
	l.wg.Add(1)
}

// Wait blocks while the limiter is not ready to proceed. Once the limiter is
// ready to go ahead, wait allows progress and notes the number of active tasks.
func (l *Limiter) Wait() {
	l.cond.L.Lock()
	for !l.ready() {
		l.cond.Wait()
	}

	l.add()
	l.cond.L.Unlock()
}

// Done tells the limiter that a task has completed and that it can allow
// another task to proceed.
func (l *Limiter) Done() {
	l.cond.L.Lock()
	l.done()
	l.cond.Signal()
	l.cond.L.Unlock()
}

// WaitDone blocks until all tasks running in the limiter call Done.
func (l *Limiter) WaitDone() {
	l.wg.Wait()
}
