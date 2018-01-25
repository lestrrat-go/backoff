.PHONY: benchmark

benchmark:
	go test -run=none -tags bench -bench . -benchmem -benchtime 20s