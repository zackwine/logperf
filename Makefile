# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=logperf
TOP=${PWD}
IMAGE_VERSION := $(shell date -u +%FT%T | sed 's/[T:-]//g')

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

# Run a simple perf config
run-cue: build
	 ./logperf/$(BINARY_NAME) -perffile ./logperf/confs/logperf-stdout-cue.yml

# Cross compilation for linux
build-linux:
	cd logperf && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../xtargets/linux/$(BINARY_NAME) -v
	echo "Built for linux: ./xtargets/linux/$(BINARY_NAME)"

get-bootstrap:
	cd web && wget -nc https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css
	cd web && wget -nc https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css.map
	cd web && wget -nc https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js
	cd web && wget -nc https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js.map
	cd web && wget -nc https://code.jquery.com/jquery-3.2.1.min.js
	cd web && wget -nc https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js
	cd web && wget -nc https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js.map

docker-build: get-bootstrap
	docker build -t logperf:$(IMAGE_VERSION) .

docker-run: docker-build
	docker run -p 8080:8080 logperf:$(IMAGE_VERSION)
