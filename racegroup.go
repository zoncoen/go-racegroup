// Package racegroup provides synchronization, and Context cancelation for
// groups of goroutines working on subtasks of a common task.
package racegroup

import (
	"context"
	"sync"
	"sync/atomic"
)

// A Group is a collection of goroutines working on subtasks.
type Group struct {
	wg     sync.WaitGroup
	cancel func()

	errHandler func(error)
	semaphore  chan struct{}
	desired    int64
	completed  int64
}

// WithContext returns a new Group and an associated Context derived from ctx.
func WithContext(ctx context.Context, opts ...Option) (*Group, context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)
	g := &Group{cancel: cancel, desired: 1}
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, nil, err
		}
	}
	return g, ctx, nil
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
// If more than or equal to desired count subtasks are completed,
// cancels the group.
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
			if atomic.AddInt64(&g.completed, 1) >= g.desired {
				g.cancel()
			}
		}
	}()
}
