GOBUILD = go build
GORUN = go run
GOBIN = ./bin
GOTEST = go test

help:
	$(GORUN) ./cmd/4337-in-a-box/main.go -h

run-app:
	$(GORUN) ./cmd/4337-in-a-box/main.go

run-test:
	@echo "Running tests..."
	$(GOTEST) -v -cover ./...

run-test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

build-source:
	$(GOBUILD) -o ./bin/4337-in-a-box ./cmd/4337-in-a-box
	@echo "Done building."
	@echo "Run \"$(GOBIN)/4337-in-a-box\" to launch 4337-in-a-box."

build-docker:
	docker build -t 4337-in-a-box:v-local .