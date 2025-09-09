.PHONY: fmt vet clean test test-it test-race test-race-cover

fmt:    ## run go fmt
	go fmt ./...

vet:    ## run go vet
	go vet ./...

clean:  ## go clean
	go clean

test:   ## run unit tests (short)
	go test -v ./... -short

test-it: ## run integration tests
	go test -v ./...

test-race: ## run tests with race detector
	go test -race -count=1 ./...

test-race-cover: ## run tests with race detector and coverage
	go test -race -count=1 -coverprofile=coverage.out ./...

test-bench: ## run benchmark tests
	go test -bench -v ./...

test-bench-mem: ## with mem and allocation
	go test -bench=. -benchmem -benchtime=2s ./...