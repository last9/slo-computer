# Makefile for SLO Computer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=slo-computer
GO111MODULE=on
GOFLAGS=-mod=vendor
DOCKER_IMAGE=last9/slo-computer

.PHONY: all build clean test run deps vendor tidy help docker docker-push docker-run

all: deps build

build:
	@echo "Building SLO Computer..."
	GO111MODULE=$(GO111MODULE) $(GOMOD) tidy
	GO111MODULE=$(GO111MODULE) $(GOBUILD) -o $(BINARY_NAME) -v

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	@echo "Running tests..."
	GO111MODULE=$(GO111MODULE) $(GOTEST) -v ./...

run: build
	@echo "Running SLO Computer..."
	./$(BINARY_NAME)

deps:
	@echo "Ensuring dependencies..."
	GO111MODULE=$(GO111MODULE) $(GOMOD) download

vendor:
	@echo "Creating vendor directory..."
	GO111MODULE=$(GO111MODULE) $(GOMOD) vendor

tidy:
	@echo "Tidying dependencies..."
	GO111MODULE=$(GO111MODULE) $(GOMOD) tidy

docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .

docker-push: docker
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):latest

docker-run: docker
	@echo "Running Docker container..."
	docker run --rm $(DOCKER_IMAGE):latest

# Example targets for common commands
example-service:
	@echo "Running service SLO example..."
	./$(BINARY_NAME) suggest --throughput=4200 --slo=99.9 --duration=720

example-cpu:
	@echo "Running CPU burst example..."
	./$(BINARY_NAME) cpu-suggest --instance=t3a.xlarge --utilization=15

example-json:
	@echo "Running service SLO example with JSON output..."
	./$(BINARY_NAME) suggest --throughput=4200 --slo=99.9 --duration=720 --output=json

# Help command
help:
	@echo "SLO Computer Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build the application after ensuring dependencies"
	@echo "  make build        Build the application"
	@echo "  make clean        Remove build artifacts"
	@echo "  make test         Run tests"
	@echo "  make run          Build and run the application"
	@echo "  make deps         Ensure dependencies are downloaded"
	@echo "  make vendor       Create vendor directory"
	@echo "  make tidy         Tidy go.mod file"
	@echo "  make docker       Build Docker image"
	@echo "  make docker-push  Push Docker image to registry"
	@echo "  make docker-run   Run Docker container"
	@echo "  make example-service  Run an example service SLO calculation"
	@echo "  make example-cpu      Run an example CPU burst calculation"
	@echo "  make example-json     Run an example with JSON output"
	@echo ""
	@echo "Environment variables:"
	@echo "  GO111MODULE       Controls Go modules behavior (default: on)" 