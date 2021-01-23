package backoff

type jitter interface {
	apply(interval float64) float64
}

func newJitter(jitterFactor float64, rng Random) jitter {
	if jitterFactor <= 0 || jitterFactor >= 1 {
		return newNopJitter()
	}
	return newRandomJitter(jitterFactor, rng)
}
