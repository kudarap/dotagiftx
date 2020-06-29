PROJECTNAME=dota2giftables

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: install build

install:
	go get ./...

run: build
	./$(PROJECTNAME)

build:
	go build -v -ldflags=" \
		-X main.tag=`git describe --tag --abbrev=0` \
		-X main.commit=`git rev-parse HEAD` \
		-X main.built=`date -u +%s`" \
		-o $(PROJECTNAME) ./cmd/$(PROJECTNAME)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags=" \
		-X main.tag=`git describe --tag --abbrev=0` \
		-X main.commit=`git rev-parse HEAD` \
		-X main.built=`date -u +%s`" \
		-o ./$(PROJECTNAME)_amd64 ./cmd/$(PROJECTNAME)

docker-build:
	docker build -t $(PROJECTNAME) .
docker-run:
	docker run -it --rm -p 8000:8000 $(PROJECTNAME)

web-build:
	cd ./web && ./build.sh .env.prod && cd ..