PROJECTNAME=dotagiftx

LDFLAGS="-X main.tag=`cat VERSION` \
		-X main.commit=`git rev-parse HEAD` \
		-X main.built=`date -u +%s`"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: install build

install:
	go get ./...

run: generate build
	./$(PROJECTNAME)

run-worker: generate build-worker
	./dxworker

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

docker-build:
	docker build -t $(PROJECTNAME) .
docker-run:
	docker run -it --rm -p 8000:8000 $(PROJECTNAME)

web-build:
	cd ./web && yarn dev && cd ..