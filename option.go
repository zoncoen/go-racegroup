package racegroup

import "errors"

// A Option is a functional option that changes the behavior of Group.
type Option func(*Group) error

// ErrorHandler returns an Option that sets the error handler.
func ErrorHandler(handler func(error)) Option {
	return func(g *Group) error {
		g.errHandler = handler
		return nil
	}
}

// Concurrency returns an Option that sets number of concurrency for goroutines.
func Concurrency(i int) Option {
	return func(g *Group) error {
		if i < 1 {
			return errors.New("concurrency option must be greater than zero")
		}
		g.semaphore = make(chan struct{}, i)
		return nil
	}
}
