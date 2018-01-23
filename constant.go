package backoff

import (
	"context"
	"time"
)

func NewConstant(delay time.Duration, options ...Option) *Constant {
	maxRetries := -1
	for _, o := range options {
		switch o.Name() {
		case optkeyMaxRetries:
			maxRetries = o.Value().(int)
		}
	}

	return &Constant{
		delay:      delay,
		maxRetries: maxRetries,
	}
}

func (p *Constant) Start(ctx context.Context) (Backoff, CancelFunc) {
	b := &constantBackoff{
		baseBackoff: newBaseBackoff(ctx, p.maxRetries),
		policy:      p,
	}

	return b, CancelFunc(b.cancelLocked)
}

func (b *constantBackoff) Next() <-chan struct{} {
	time.AfterFunc(b.policy.delay, b.fire)
	return b.next
}
