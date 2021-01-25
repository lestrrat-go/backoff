module github.com/lestrrat-go/backoff/bench

go 1.16

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/lestrrat-go/backoff/v2 v2.0.0-00010101000000-000000000000
)

replace github.com/lestrrat-go/backoff/v2 => ../
