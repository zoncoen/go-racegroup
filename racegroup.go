// Package racegroup provides synchronization, and Context cancelation for
// groups of goroutines working on subtasks of a common task.
package racegroup

import (
	"context"
	"sync"
)

// A Group is a collection of goroutines working on subtasks.
type Group struct {
	wg     sync.WaitGroup
	cancel func()

	errHandler func(error)
	semaphore  chan struct{}
}

// WithContext returns a new Group and an associated Context derived from ctx.
func WithContext(ctx context.Context, opts ...Option) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	g := &Group{cancel: cancel}
	for _, opt := range opts {
		opt(g)
	}
	return g, ctx
}

// Wait blocks until all function calls from the Go method have returned.
func (g *Group) Wait() {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
}

// Go calls the given function in a new goroutine.
//
// The first call to return a nil error cancels the group.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	if g.semaphore != nil {
		g.semaphore <- struct{}{}
	}

	go func() {
		defer g.wg.Done()
		defer func() {
			if g.semaphore != nil {
				<-g.semaphore
			}
		}()

		if err := f(); err != nil {
			if g.errHandler != nil {
				g.errHandler(err)
			}
		} else {
			g.cancel()
		}
	}()
}
