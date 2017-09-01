package racegroup_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	racegroup "github.com/zoncoen/go-racegroup"
)

func wait(ctx context.Context, d time.Duration) func() error {
	return func() error {
		select {
		case <-time.After(d):
			fmt.Printf("wait %s\n", d)
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}
}

func errFunc() func() error {
	return func() error {
		return errors.New("error occurred")
	}
}

func errPrinter(err error) {
	fmt.Println(err)
}

func ExampleGroup() {
	g, ctx, _ := racegroup.WithContext(context.Background())
	g.Go(wait(ctx, 2*time.Second))
	g.Go(wait(ctx, 1*time.Second))
	g.Wait()

	// Output:
	// wait 1s
}

func ExampleErrorHandler() {
	g, ctx, _ := racegroup.WithContext(context.Background(), racegroup.ErrorHandler(errPrinter))
	g.Go(wait(ctx, 1*time.Second))
	g.Go(errFunc())
	g.Wait()

	// Output:
	// error occurred
	// wait 1s
}

func ExampleConcurrency() {
	g, ctx, _ := racegroup.WithContext(context.Background(), racegroup.Concurrency(2))
	g.Go(wait(ctx, 3*time.Second))
	g.Go(wait(ctx, 2*time.Second))
	g.Go(wait(ctx, 1*time.Second))
	g.Wait()

	// Output:
	// wait 2s
}
