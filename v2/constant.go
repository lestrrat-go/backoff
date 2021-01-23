package backoff

import (
	"context"
	"math/rand"
	"time"
)

type ConstantInterval struct {
	interval     time.Duration
	jitterFactor float64
	rng          Random
}

func NewConstantInterval(options ...ConstantOption) *ConstantInterval {
	jitterFactor := 0.0
	interval := time.Minute
	var rng Random

	for _, option := range options {
		switch option.Ident() {
		case identInterval{}:
			interval = option.Value().(time.Duration)
		case identJitterFactor{}:
			jitterFactor = option.Value().(float64)
		case identRNG{}:
			rng = option.Value().(Random)
		}
	}

	if jitterFactor <= 0 || jitterFactor >= 1 {
		jitterFactor = 0
	}

	if rng == nil {
		// if we have a jitter factor, and no RNG is provided, create one.
		// This is definitely not "secure", but well, if you care enough,
		// you would provide one
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return &ConstantInterval{
		interval:     interval,
		jitterFactor: jitterFactor,
		rng:          rng,
	}
}

func (g *ConstantInterval) Next() time.Duration {
	if factor := g.jitterFactor; factor > 0 {
		interval := float64(g.interval)
		jitterDelta := interval * factor
		jitterMin := interval - jitterDelta
		jitterMax := interval + jitterDelta

		return time.Duration(jitterMin + g.rng.Float64()*(jitterMax-jitterMin+1))
	}

	return g.interval
}

type ConstantPolicy struct {
	cOptions  []ControllerOption
	igOptions []ConstantOption
}

func NewConstantPolicy(options ...Option) *ConstantPolicy {
	var cOptions []ControllerOption
	var igOptions []ConstantOption

	for _, option := range options {
		switch opt := option.(type) {
		case ControllerOption:
			cOptions = append(cOptions, opt)
		default:
			igOptions = append(igOptions, opt.(ConstantOption))
		}
	}

	return &ConstantPolicy{
		cOptions:  cOptions,
		igOptions: igOptions,
	}
}

func (p *ConstantPolicy) Start(ctx context.Context) Controller {
	ig := NewConstantInterval(p.igOptions...)
	return newController(ctx, ig, p.cOptions...)
}
