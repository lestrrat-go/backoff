package backoff

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

const (
	optkeyFactor         = "factor"
	optkeyInterval       = "interval"
	optkeyJitterFactor   = "jitter-factor"
	optkeyMaxElapsedTime = "max-elapsed-time"
	optkeyMaxInterval    = "max-interval"
	optkeyMaxRetries     = "max-retries"
	optkeyThreshold      = "threshold"
)

const (
	defaultInterval     = 500 * time.Millisecond
	defaultJitterFactor = 0.5
	defaultMaxInterval  = float64(2 * time.Minute)
	defaultMaxRetries   = 10
	defaultThreshold    = 15 * time.Minute
)

type CancelFunc func()

type Policy interface {
	Start(context.Context) (Backoff, CancelFunc)
}

type Backoff interface {
	Done() <-chan struct{}
	Next() <-chan struct{}
}

type Constant struct {
	delay          time.Duration
	maxElapsedTime time.Duration
	maxRetries     int
}

type Option interface {
	Name() string
	Value() interface{}
}

type baseBackoff struct {
	current        interface{}
	callCount      int
	cancelFunc     context.CancelFunc
	ctx            context.Context
	maxElapsedTime time.Duration
	maxRetries     int
	mu             sync.RWMutex
	startTime      time.Time
	next           chan struct{}
}

type constantBackoff struct {
	*baseBackoff
	policy *Constant
}

// Exponential implements an exponential backoff policy.
type Exponential struct {
	factor         float64
	interval       time.Duration
	jitterFactor   float64
	maxElapsedTime time.Duration
	maxInterval    float64
	maxRetries     int
	random         *rand.Rand
	threshold      time.Duration // max backoff
}

type exponentialBackoff struct {
	*baseBackoff
	policy  *Exponential
	attempt float64
}
