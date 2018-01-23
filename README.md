# go-backoff

Backoff algorithm and helpers for Go

# SYNOPSIS

```go

import "github.com/lestrrat/go-backoff"

func Func(arg Foo) (Result, error) { ... }

var policy = backoff.NewExponential(
  backoff.WithInterval(500*time.Millisecond), // base interval
  backoff.WithJitterFactor(0.05), // 5% jitter
  backoff.WithMaxRetries(25),
)
func RetryFunc(arg Foo) (Result, error) {
  b := policy.Start(context.Background())
  for {
    res, err := Func(arg)
    if err == nil {
      return res, nil
    }

    select {
    case <-b.Done():
      return nil, errors.New(`tried very hard, but no luck`)
    case <-b.Next():
      // at this point we can fire the next call, so
      // just continue with the loop
    }
  }

  return nil, errors.New(`unreachable`)
}
```

# DESCRIPTION

This library is an implementation of backoff algorithm for retrying operations
in an idiomatic Go way. It respects `context.Context` natively, and the critical
notifications are done through channel operations, allowing you greater
flexibility in how you wrap your operations

# PRIOR ARTS

## [github.com/cenkalti/backoff](https://github.com/cenkalti/backoff) 

This library is featureful, but one thing that gets to me is the fact that
it essentially forces you to create a closure over the operation to be retried.

## [github.com/jpillora/backoff](https://github.com/jpillora/backoff)

This library is a very simple implementation of calculating backoff durations.
I wanted it to let me know when to stop too, so it was missing a few things.