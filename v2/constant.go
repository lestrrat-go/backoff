package backoff

import (
	"context"
	"time"
)

type ConstantInterval struct {
	interval time.Duration
}

func NewConstantInterval(options ...ConstantOption) *ConstantInterval {
	var interval time.Duration = time.Minute
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
