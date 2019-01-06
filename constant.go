package backoff

import (
	"context"
	"time"
)

func NewConstant(delay time.Duration, options ...Option) *Constant {
	maxRetries := defaultMaxRetries
	var maxElapsedTime time.Duration
	for _, o := range options {
		switch o.Name() {
		case optkeyMaxElapsedTime:
			maxElapsedTime = o.Value().(time.Duration)
		case optkeyMaxRetries:
			maxRetries = o.Value().(int)
		}
	}

	return &Constant{
		delay:          delay,
		maxElapsedTime: maxElapsedTime,
		maxRetries:     maxRetries,
	}
}

func (p *Constant) Start(ctx context.Context) (Backoff, CancelFunc) {
	b := &constantBackoff{
		baseBackoff: newBaseBackoff(ctx, p.maxRetries, p.maxElapsedTime),
		policy:      p,
	}
	b.baseBackoff.Start(ctx)

	b.mu.Lock()
	b.current = 1 // record that we've already queued the first fake event
	go b.fire()   // the first call
	b.mu.Unlock()

	return b, CancelFunc(b.cancelLocked)
}

func (b *constantBackoff) Next() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Only queue a request to fire if the previous request
	// has already been processed
	if b.current == nil {
		b.current = time.AfterFunc(b.policy.delay, b.fire)
	}
	return b.next
}
