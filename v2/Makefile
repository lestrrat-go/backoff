.PHONY: realclean cover viewcover test lint

realclean:
	rm coverage.out

test:
	cd examples && go test -v -race && cd .. && go test -v -race ./...

cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.out.tmp ./...
	@# This is NOT cheating. tools to generate code don't need to be
	@# included in the final result
	@cat coverage.out.tmp | grep -v "internal/cmd" > coverage.out
	@rm coverage.out.tmp

smoke:
	go test -race -short ./...

viewcover:
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...
