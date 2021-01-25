package backoff_test

import (
	"context"
	"errors"
	"time"

	backoff "github.com/lestrrat-go/backoff/v2"
)

func ExampleConstant() {
	p := backoff.Constant(backoff.WithInterval(time.Second))

	flakyFunc := func(a int) (int, error) {
		// silly function that only succeeds if the current call count is
		// divisible by either 3 or 5 but not both
		switch {
		case a%15 == 0:
			return 0, errors.New(`invalid`)
		case a%3 == 0 || a%5 == 0:
			return a, nil
		}
		return 0, errors.New(`invalid`)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	retryFunc := func(v int) (int, error) {
		b := p.Start(ctx)
		for backoff.Continue(b) {
			x, err := flakyFunc(v)
			if err == nil {
				return x, nil
			}
		}
		return 0, errors.New(`failed to get value`)
	}

	retryFunc(15)
}

func ExampleExponential() {
	p := backoff.Exponential(
		backoff.WithMinInterval(time.Second),
		backoff.WithMaxInterval(time.Minute),
		backoff.WithJitterFactor(0.05),
	)

	flakyFunc := func(a int) (int, error) {
		// silly function that only succeeds if the current call count is
		// divisible by either 3 or 5 but not both
		switch {
		case a%15 == 0:
			return 0, errors.New(`invalid`)
		case a%3 == 0 || a%5 == 0:
			return a, nil
		}
		return 0, errors.New(`invalid`)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	retryFunc := func(v int) (int, error) {
		b := p.Start(ctx)
		for backoff.Continue(b) {
			x, err := flakyFunc(v)
			if err == nil {
				return x, nil
			}
		}
		return 0, errors.New(`failed to get value`)
	}

	retryFunc(15)
}
