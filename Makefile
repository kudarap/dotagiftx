PROJECTNAME=dota2gitables

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: install build

install:
	go get ./...

run:
	make build && ./api

build:
	go build -v -ldflags=" \
		-X main.tag=`git describe --tag --abbrev=0` \
		-X main.commit=`git rev-parse HEAD` \
		-X main.built=`date -u +%s`" \
		-o api ./cmd/api
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./api_unix ./cmd/api

docker-build:
	docker build -t $(PROJECTNAME) .
docker-run:
	docker run -it --rm -p 8000:8000 $(PROJECTNAME)

web-build:
	cd ./web && ./build.sh .env.prod && cd ..