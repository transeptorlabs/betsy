GOBUILD = go build
GORUN = go run
GOBIN = ./bin
GOTEST = go test

help:
	$(GORUN) ./cmd/betsy/main.go -h

run-cli:
	$(GORUN) ./cmd/betsy/main.go

run-cli-dev:
	$(GORUN) ./cmd/betsy/main.go --log DEBUG --debug

run-test:
	@echo "Running tests..."
	$(GOTEST) -v -cover ./...

run-test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

betsy:
	$(GOBUILD) -o ./bin/betsy ./cmd/betsy
	@echo "Done building."
	@echo "Run \"$(GOBIN)/betsy\" to launch betsy."

gen-contract-binding-aa:
	@echo "Generating contract bindings..."
	chmod +x ./scripts/gen-contracts-binding-aa.sh
	./scripts/gen-contracts-binding-aa.sh