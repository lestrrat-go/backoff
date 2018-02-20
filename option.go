package backoff

import "time"

type option struct {
	name  string
	value interface{}
}

func (o option) Name() string       { return o.name }
func (o option) Value() interface{} { return o.value }

func WithFactor(v float64) Option {
	return &option{
		name:  optkeyFactor,
		value: v,
	}
}

func WithInterval(v time.Duration) Option {
	return &option{
		name:  optkeyInterval,
		value: v,
	}
}

func WithJitterFactor(v float64) Option {
	return &option{
		name:  optkeyJitterFactor,
		value: v,
	}
}

// WithMaxInterval specifies the maximum interval between retries, and is
// currently only applicable to exponential backoffs.
//
// By default this is capped at 2 minutes. If you would like to change this
// value, you must explicitly specify it through this option.
//
// If a value of 0 is specified, then there is no limit, and the backoff
// interval will keep increasing.
func WithMaxInterval(v time.Duration) Option {
	return &option{
		name:  optkeyMaxInterval,
		value: float64(v),
	}
}

// WithMaxRetries specifies the maximum number of attempts that can be made
// by the backoff policies. By default each policy tries up to 10 times.
//
// If you would like to retry forever, specify "0" and pass to the constructor
// of each policy.
func WithMaxRetries(v int) Option {
	return &option{
		name:  optkeyMaxRetries,
		value: v,
	}
}

// WithMaxElapsedTime specifies the maximum amount of accumulative time that
// the backoff is allowed to wait before it is considered failed.
func WithMaxElapsedTime(v time.Duration) Option{
	return &option{
		name: optkeyMaxElapsedTime,
		value: v,
	}
}
