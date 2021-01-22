# go-backoff

Backoff algorithm and helpers for Go

This path (github.com/lestrrat-go/backoff) points to v1.
** PLEASE USE github.com/lestrrat-go/backoff/v ** for currently supported version


# SYNOPSIS

```go

import "github.com/lestrrat-go/backoff"

func Func(arg Foo) (Result, error) { ... }

var policy = backoff.Exponential(
  backoff.WithMinInterval(500*time.Millisecond), // base interval
  backoff.WithJitterFactor(0.05), // 5% jitter
  backoff.WithMaxRetries(25), // If not specified, default number of retries is 10
)
func RetryFunc(arg Foo) (Result, error) {
  b, cancel := policy.Start(context.Background())
  defer cancel()

  for backoff.Continue(b) {
    res, err := Func(arg)
    if err == nil {
      return res, nil
    }
  }
  return nil, errors.New(`tried very hard, but no luck`)
}
```

# DESCRIPTION

This library is an implementation of backoff algorithm for retrying operations
in an idiomatic Go way. It respects `context.Context` natively, and the critical
notifications are done through *channel operations*, allowing you to write code 
that is both more explicit and flexibile.

It also exports a utility function `Retry`, for simple operations.

For a longer discussion, [please read this article](https://medium.com/@lestrrat/yak-shaving-with-backoff-libraries-in-go-80240f0aa30c)

# PRIOR ARTS

## [github.com/cenkalti/backoff](https://github.com/cenkalti/backoff) 

This library is featureful, but one thing that gets to me is the fact that
it essentially forces you to create a closure over the operation to be retried.

## [github.com/jpillora/backoff](https://github.com/jpillora/backoff)

This library is a very simple implementation of calculating backoff durations.
I wanted it to let me know when to stop too, so it was missing a few things.
