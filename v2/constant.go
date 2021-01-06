package backoff

import (
	"context"
	"time"
)

type ConstantInterval struct {
	interval time.Duration
}

func NewConstantInterval(options ...Option) *ConstantInterval {
	var interval time.Duration = 15 * time.Minute
	for _, option := range options {
		switch option.Ident() {
		case identInterval{}:
			interval = option.Value().(time.Duration)
		}
	}

	return &ConstantInterval{
		interval: interval,
	}
}

func (g *ConstantInterval) Next() time.Duration {
	return g.interval
}

type ConstantPolicy struct {
	cOptions []Option
	igOptions []Option
}

func NewConstantPolicy(options ...Option) *ConstantPolicy {
	var cOptions []Option
	var igOptions []Option

	for _, option := range options {
		switch option.Ident() {
		case identInterval{}:
			igOptions = append(igOptions, option)
		default:
			cOptions = append(cOptions, option)
		}
	}

	return &ConstantPolicy{
		cOptions: cOptions,
		igOptions: igOptions,
	}
}

func (p *ConstantPolicy) Start(ctx context.Context) Controller {
	ig := NewConstantInterval(p.igOptions...)
	return newController(ctx, ig, p.cOptions...)
}
