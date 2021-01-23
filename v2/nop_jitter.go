package backoff

type nopJitter struct{}

func newNopJitter() *nopJitter {
	return &nopJitter{}
}

func (j *nopJitter) apply(interval float64) float64 {
	return interval
}
