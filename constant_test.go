package backoff_test

import (
	"testing"
	"time"

	backoff "github.com/lestrrat-go/backoff"
)

func TestConstantInterface(t *testing.T) {
	var b backoff.Policy = backoff.NewConstant(time.Second)
	_ = b
}
