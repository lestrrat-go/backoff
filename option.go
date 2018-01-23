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

func WithMaxRetries(v int) Option {
	return &option{
		name:  optkeyMaxRetries,
		value: v,
	}
}
