package racegroup

// A Option is a functional option that changes the behavior of Group.
type Option func(*Group)

// ErrorHandler returns an Option that sets the error handler.
func ErrorHandler(handler func(error)) Option {
	return func(g *Group) {
		g.errHandler = handler
	}
}

// Concurrency returns an Option that sets number of concurrency for goroutines.
func Concurrency(i int) Option {
	return func(g *Group) {
		g.semaphore = make(chan struct{}, i)
	}
}
