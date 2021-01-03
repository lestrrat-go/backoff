package backoff

import (
	"time"

	"github.com/lestrrat-go/option"
)

type identInterval struct{}
type identJitterFactor struct{}
type identMaxInterval struct{}
type identMaxRetries struct{}
type identMinInterval struct{}
type identMultiplier struct{}
type identRNG struct{}

// WithInterval specifies the constant interval used in ConstantPolicy and
// ConstantInterval.
func WithInterval(v time.Duration) Option {
	return option.New(identInterval{}, v)
}

// WithMaxRetries specifies the maximum number of attempts that can be made
// by the backoff policies. By default each policy tries up to 10 times.
//
// If you would like to retry forever, specify "0" and pass to the constructor
// of each policy.
//
// This option can be passed to all policy constructors except for NullPolicy
func WithMaxRetries(v int) Option {
	return option.New(identMaxRetries{}, v)
}

// WithMaxInterval specifies the maximum duration used ax exponential backoff
// The default value is 1 minute.
//
// This option can be passed to ExponentialPolicy constructor
func WithMaxInterval(v time.Duration) Option {
	return option.New(identMaxInterval{}, v)
}

// WithMinInterval specifies the minimum duration used in exponential backoff.
// The default value is 500ms.
//
// This option can be passed to ExponentialPolicy constructor
func WithMinInterval(v time.Duration) Option {
	return option.New(identMinInterval{}, v)
}

// WithMultiplier specifies the factor in which the backoff intervals are
// increased. By default this value is set to 1.5, which means that for
// every iteration a 50% increase in the interval for every iteration
// (up to the value specified by WithMaxInterval). this value must be greater
// than 1.0. If the value is less than equal to 1.0, the default value
// of 1.5 is used.
//
// This option can be passed to ExponentialPolicy constructor
func WithMultiplier(v float64) Option {
	return option.New(identJitterFactor{}, v)
}

// WithJitterFactor enables some randomness (jittering) in the computation of
// the backoff intervals. This value must be between 0.0 < v < 1.0. If a
// value outside of this range is specified, the value will be silently
// ignored and jittering is disabled.
//
// This option can be passed to ExponentialPolicy constructor
func WithJitterFactor(v float64) Option {
	return option.New(identJitterFactor{}, v)
}

// WithRNG specifies the random number generator used for jittering.
// If not provided one will be created, but if you want a truly random
// jittering, make sure to provide one that you explicitly initialized
func WithRNG(v Random) Option {
	return option.New(identRNG{}, v)
}
