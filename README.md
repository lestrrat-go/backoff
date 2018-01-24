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
  b, cancel := policy.Start(context.Background())
  defer cancel()

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

# DUMB BENCHMARK

```
go test -run=none -bench . -tags bench -benchmem -benchtime 20s
goos: darwin
goarch: amd64
pkg: github.com/lestrrat/go-backoff
Benchmark/cenkalti-4         	       5	6390580465 ns/op	    1171 B/op	      24 allocs/op
Benchmark/lestrrat-4         	       5	5077630205 ns/op	    1059 B/op	      21 allocs/op
PASS
ok  	github.com/lestrrat/go-backoff	102.885s
```