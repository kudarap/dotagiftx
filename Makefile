# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

server_bin=dxserver
worker_bin=dxworker
build_flags="-X main.tag=`cat VERSION` -X main.commit=`git rev-parse HEAD` -X main.built=`date -u +%s`"

all: test fmt build build-linux build-worker build-worker-linux

install:
	go get ./...

run: build
	./$(server_bin)

run-worker: build-worker
	./$(worker_bin)

test: lint
	go test -v ./
	go test -v ./http/...
	go test -v ./steam/...
	go test -v ./phantasm/...
	go test -v ./verify/...

fmt: generate
	gofmt -s -l -e -w .

lint:
	golangci-lint run -v

generate:
	go generate .

build:
	GOEXPERIMENT=greenteagc GOEXPERIMENT=jsonv2 go build -v -ldflags=$(build_flags) -o $(server_bin) ./cmd/$(server_bin)
build-worker:
	GOEXPERIMENT=greenteagc GOEXPERIMENT=jsonv2 go build -v -ldflags=$(build_flags) -o $(worker_bin) ./cmd/$(worker_bin)
build-linux:
	GOEXPERIMENT=greenteagc GOEXPERIMENT=jsonv2 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
		-ldflags=$(build_flags) -o ./$(server_bin)_amd64 ./cmd/$(server_bin)
build-worker-linux:
	GOEXPERIMENT=greenteagc GOEXPERIMENT=jsonv2 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
		-ldflags=$(build_flags) -o $(worker_bin)_amd64 ./cmd/$(worker_bin)

docker-build:
	docker build -t $(server_bin) .
docker-run:
	docker run -it --rm -p 8000:8000 $(server_bin)

web-build:
	cd ./web && yarn dev && cd ..

local:
	docker-compose up