package backoff

import (
	"context"
	"math/rand"
	"time"
)

type ExponentialInterval struct {
	current      float64
	jitterFactor float64
	maxInterval  float64
	minInterval  float64
	multiplier   float64
	rng          Random
}

const (
	defaultMaxInterval = float64(time.Minute)
	defaultMinInterval = float64(500 * time.Millisecond)
	defaultMultiplier  = 1.5
)

func NewExponentialInterval(options ...ExponentialOption) *ExponentialInterval {
	jitterFactor := 0.0
	maxInterval := defaultMaxInterval
	minInterval := defaultMinInterval
	multiplier := defaultMultiplier
	var rng Random

	for _, option := range options {
		switch option.Ident() {
		case identJitterFactor{}:
			jitterFactor = option.Value().(float64)
		case identMaxInterval{}:
			maxInterval = float64(option.Value().(time.Duration))
		case identMinInterval{}:
			minInterval = float64(option.Value().(time.Duration))
		case identMultiplier{}:
			multiplier = option.Value().(float64)
		case identRNG{}:
			rng = option.Value().(Random)
		}
	}

	if minInterval > maxInterval {
		minInterval = maxInterval
	}
	if multiplier <= 1 {
		multiplier = defaultMultiplier
	}
	if jitterFactor <= 0 || jitterFactor >= 1 {
		jitterFactor = 0
	}
	if jitterFactor > 0 && rng == nil {
		// if we have a jitter factor, and no RNG is provided, create one.
		// This is definitely not "secure", but well, if you care enough,
		// you would provide one
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return &ExponentialInterval{
		jitterFactor: jitterFactor,
		maxInterval:  maxInterval,
		minInterval:  minInterval,
		multiplier:   multiplier,
		rng:          rng,
	}
}

func (g *ExponentialInterval) Next() time.Duration {
	var next float64
	if g.current == 0 {
		next = g.minInterval
	} else {
		next = g.current * g.multiplier
	}

	if next > g.maxInterval {
		next = g.maxInterval
	}
	if next < g.minInterval {
		next = g.minInterval
	}

	// Apply jitter *AFTER* we calculate the base interval
	if factor := g.jitterFactor; factor > 0 {
		jitterDelta := next * factor
		jitterMin := next - jitterDelta
		jitterMax := next + jitterDelta

		next = jitterMin + g.rng.Float64()*(jitterMax-jitterMin+1)
	}

	g.current = next
	return time.Duration(next)
}

type ExponentialPolicy struct {
	cOptions  []ControllerOption
	igOptions []ExponentialOption
}

func NewExponentialPolicy(options ...ExponentialOption) *ExponentialPolicy {
	var cOptions []ControllerOption
	var igOptions []ExponentialOption

	for _, option := range options {
		switch opt := option.(type) {
		case ControllerOption:
			cOptions = append(cOptions, opt)
		default:
			igOptions = append(igOptions, opt)
		}
	}

	return &ExponentialPolicy{
		cOptions:  cOptions,
		igOptions: igOptions,
	}
}

func (p *ExponentialPolicy) Start(ctx context.Context) Controller {
	ig := NewExponentialInterval(p.igOptions...)
	return newController(ctx, ig, p.cOptions...)
}

