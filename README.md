# go-backoff

Backoff algorithm and helpers for Go

# SYNOPSIS

```go

func Func(arg Foo) (Result, error) { ... }

var policy = backoff.NewExponential()
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

* Works well with context.Context
* Does not require closures

# PRIOR ARTS

## [https://github.com/cenkalti/backoff](github.com/cenkalti/backoff) 

This library is featureful, but one thing that gets to me is the fact that
it essentially forces you to create a closure over the operation to be retried.

## [https://github.com/jpillora/backoff](github.com/jpillora/backoff)

This library is a very simple implementation of calculating backoff durations.
I wanted it to let me know when to stop too, so it was missing a few things.