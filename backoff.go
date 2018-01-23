package backoff

import "context"

func newBaseBackoff(ctx context.Context, maxRetries int) *baseBackoff {
	backoffCtx, cancel := context.WithCancel(ctx)
	return &baseBackoff{
		cancelFunc: cancel,
		ctx:        backoffCtx,
		maxRetries: maxRetries,
		next:       make(chan struct{}),
	}
}

func (b *baseBackoff) Done() <-chan struct{} {
	return b.ctx.Done()
}

func (b *baseBackoff) cancelLocked() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.cancel()
}

// note: caller must lock
func (b *baseBackoff) cancel() {
	close(b.next)
	b.next = nil
	b.cancelFunc()
}

func (b *baseBackoff) fire() {
	select {
	case <-b.ctx.Done():
		return
	default:
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.next == nil {
		return
	}

	b.next <- struct{}{}
	if b.maxRetries > 0 {
		if b.maxRetries <= b.callCount {
			b.cancel()
		} else {
			b.callCount++
		}
	}
}
