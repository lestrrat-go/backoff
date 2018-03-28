package backoff_test

import (
	"context"
	"errors"
	"log"
	"time"

	backoff "github.com/lestrrat-go/backoff"
)

func Example() {
	p := backoff.NewConstant(time.Second)

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
		b, cancel := p.Start(ctx)
		defer cancel()

		for {
			x, err := flakyFunc(v)
			if err == nil {
				return x, nil
			}

			select {
			case <-b.Done():
				return 0, errors.New(`....`)
			case <-b.Next():
				// no op, go to next
			}
		}
	}

	retryFunc(15)
}

func ExampleContinue() {
	p := backoff.NewConstant(time.Second)

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
		b, cancel := p.Start(ctx)
		defer cancel()

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

func ExampleRetry() {
	count := 0
	e := backoff.ExecuteFunc(func(_ context.Context) error {
		// This is a silly example that succeeds on every 10th try
		count++
		if count%10 == 0 {
			return nil
		}
		return errors.New(`dummy`)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := backoff.NewExponential()
	if err := backoff.Retry(ctx, p, e); err != nil {
		log.Printf("failed to call function after repeated tries")
	}
}
