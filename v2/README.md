# backoff

Idiomatic backoff for Go

This library is an implementation of backoff algorithm for retrying operations
in an idiomatic Go way. It respects `context.Context` natively, and the critical
notifications are done through *channel operations*, allowing you to write code 
that is both more explicit and flexibile.

For a longer discussion, [please read this article](https://medium.com/@lestrrat/yak-shaving-with-backoff-libraries-in-go-80240f0aa30c)

# IMPORT

```go
import "github.com/lesrtrrat-go/backoff/v2"
```

# SYNOPSIS

```go
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
```
