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
