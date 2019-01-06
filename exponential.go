package backoff

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)

func NewExponential(options ...Option) *Exponential {
	interval := defaultInterval
	jitterFactor := defaultJitterFactor
	maxInterval := defaultMaxInterval
	maxRetries := defaultMaxRetries
	threshold := defaultThreshold
	factor := float64(2)
	var maxElapsedTime time.Duration
	for _, o := range options {
		switch o.Name() {
		case optkeyFactor:
			factor = o.Value().(float64)
		case optkeyInterval:
			interval = o.Value().(time.Duration)
		case optkeyJitterFactor:
			jitterFactor = o.Value().(float64)
		case optkeyMaxElapsedTime:
			maxElapsedTime = o.Value().(time.Duration)
		case optkeyMaxInterval:
			maxInterval = o.Value().(float64)
		case optkeyMaxRetries:
			maxRetries = o.Value().(int)
		case optkeyThreshold:
			threshold = o.Value().(time.Duration)
		}
	}

	return &Exponential{
		factor:         factor,
		interval:       interval,
		jitterFactor:   jitterFactor,
		maxElapsedTime: maxElapsedTime,
		maxInterval:    maxInterval,
		maxRetries:     maxRetries,
		random:         rand.New(rand.NewSource(time.Now().UnixNano())),
		threshold:      threshold,
	}
}

var exponentialBackoffPool = sync.Pool{
	New: func() interface{} {
		return &exponentialBackoff{}
	},
}

func getExponentialBackoff() *exponentialBackoff {
	return exponentialBackoffPool.Get().(*exponentialBackoff)
}

func releaseExponentialBackoff(b *exponentialBackoff) {
	b.baseBackoff = nil
	b.policy = nil
	exponentialBackoffPool.Put(b)
}

func (p *Exponential) Start(ctx context.Context) (Backoff, CancelFunc) {
	b := getExponentialBackoff()
	b.baseBackoff = newBaseBackoff(ctx, p.maxRetries, p.maxElapsedTime)
	b.policy = p
	b.attempt = 0
	b.baseBackoff.Start(ctx)

	b.mu.Lock()
	b.current = 1 // record that we've already queued the first fake event
	go b.fire()   // the first call
	b.mu.Unlock()

	return b, CancelFunc(func() {
		b.cancelLocked()
		releaseExponentialBackoff(b)
	})
}

func (b *exponentialBackoff) Next() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.current == nil {
		d := b.delayForAttempt(b.attempt)
		b.attempt++
		b.current = time.AfterFunc(d, b.fire)
	}

	return b.next
}

func (b *exponentialBackoff) delayForAttempt(attempt float64) time.Duration {
	minf := float64(b.policy.interval)
	durf := minf * math.Pow(b.policy.factor, attempt)
	if b.policy.jitterFactor > 0 {
		jitterDelta := durf * b.policy.jitterFactor
		jitteredMin := durf - jitterDelta
		jitteredMax := durf + jitterDelta

		durf = jitteredMin + b.policy.random.Float64()*(jitteredMax-jitteredMin+1)
	}

	if maxf := b.policy.maxInterval; maxf > 0 && durf > maxf {
		durf = maxf
	}

	dur := time.Duration(durf)
	return dur
}
