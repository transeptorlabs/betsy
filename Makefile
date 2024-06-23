GOBUILD = go build
GORUN = go run
GOBIN = ./bin
GOTEST = go test

eth:
	chmod +x ./scripts/eth-node.sh
	./scripts/eth-node.sh

help:
	$(GORUN) ./cmd/betsy/main.go -h

run-app:
	$(GORUN) ./cmd/betsy/main.go

run-test:
	@echo "Running tests..."
	$(GOTEST) -v -cover ./...

run-test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

build-source:
	$(GOBUILD) -o ./bin/betsy ./cmd/betsy
	@echo "Done building."
	@echo "Run \"$(GOBIN)/betsy\" to launch betsy."

build-docker:
	docker build -t betsy:v-local .