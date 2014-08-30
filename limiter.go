/*
Package limiter provides a singular type that allows you to control concurrent
execution by specifying a maximum number of tasks.

It allows you to then block until tasks complete before spawning new ones, and
also block on all tasks until they have completed.
*/
package limiter

import "sync"

type Limiter struct {
	cond        *sync.Cond
	wg          *sync.WaitGroup
	active, max int
}

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

func (l *Limiter) Wait() {
	l.cond.L.Lock()
	for !l.ready() {
		l.cond.Wait()
	}

	l.add()
	l.cond.L.Unlock()
}

func (l *Limiter) Done() {
	l.cond.L.Lock()
	l.done()
	l.cond.Signal()
	l.cond.L.Unlock()
}

func (l *Limiter) WaitDone() {
	l.wg.Wait()
}
