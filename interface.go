package backoff

import (
	"context"
	"time"

	"github.com/lestrrat-go/option"
)

type Option = option.Interface

type Controller interface {
	Done() <-chan struct{}
	Next() <-chan struct{}
}

type IntervalGenerator interface {
	Next() time.Duration
}

// Policy is an interface for the backoff policies that this package
// implements. Users must create a controller object from this
// policy to actually do anything with it
type Policy interface {
	// Start creates a new Controller object for the backoff.
	//
	// The Controller starts a new gouroutine that will control when
	// the backoff events are fired. This goroutine may or may not
	// run forever depending on your settings.
	//
	// Therefore the caller is expected to properly terminate the
	// `contexts.Context` object that was provided after the
	// desired backoff operation has been performed. Otherwise
	// you may end up with leaked goroutines.
	Start(context.Context) Controller
}

type Random interface {
	Float64() float64
}
