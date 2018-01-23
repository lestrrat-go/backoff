package backoff_test

import (
	"context"
	"errors"
	"time"

	backoff "github.com/lestrrat/go-backoff"
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
