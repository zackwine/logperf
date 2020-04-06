# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=logperf
TOP=${PWD}

all: test build
config:
	go mod edit -replace github.com/zackwine/logperf=$(TOP)
build:
	cd logperf && $(GOBUILD) -o $(BINARY_NAME) -v
	echo "Built for host: ./logperf/$(BINARY_NAME)"
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f ./logperf/$(BINARY_NAME) ./xtargets/linux/$(BINARY_NAME)

# Run a simple perf config
run: build
	 ./logperf/$(BINARY_NAME) -perffile ./logperf/confs/logperf-stdout-short.yml

# Cross compilation for linux
build-linux:
	cd logperf && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../xtargets/linux/$(BINARY_NAME) -v
	echo "Built for linux: ./xtargets/linux/$(BINARY_NAME)"
