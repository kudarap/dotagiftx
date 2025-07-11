PROJECTNAME=dotagiftx

LDFLAGS="-X main.tag=`cat VERSION` \
		-X main.commit=`git rev-parse HEAD` \
		-X main.built=`date -u +%s`"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: test build build-linux build-worker build-worker-linux

install:
	go get ./...

run: test build
	./$(PROJECTNAME)

run-worker: test build-worker
	./dxworker

test: generate fmt
	go test -v ./
	go test -v ./http/...
	go test -v ./steam/...
	go test -v ./phantasm/...
	go test -v ./verify/...

build:
	go build -v -ldflags=$(LDFLAGS) -o $(PROJECTNAME) ./cmd/$(PROJECTNAME)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags=$(LDFLAGS) \
		-o ./$(PROJECTNAME)_amd64 ./cmd/$(PROJECTNAME)
build-worker:
	go build -v -ldflags=$(LDFLAGS) -o dxworker ./cmd/dxworker
build-worker-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags=$(LDFLAGS) \
	    -o dxworker_amd64 ./cmd/dxworker

generate:
	go generate .

fmt:
	gofmt -s -l -e -w .

docker-build:
	docker build -t $(PROJECTNAME) .
docker-run:
	docker run -it --rm -p 8000:8000 $(PROJECTNAME)

web-build:
	cd ./web && yarn dev && cd ..

local:
	docker-compose up