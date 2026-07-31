# build stage
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git make curl

ENV GOCACHE=/gobuildcache
ENV GOMODCACHE=/gomodcache

WORKDIR /code

# download and cache go dependencies
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/gomodcache go mod download

# then copy source code as the last step
COPY . .

RUN --mount=type=cache,target=/gomodcache --mount=type=cache,target=/gobuildcache \
    go build \
    -ldflags="-X main.tag=`cat VERSION` -X main.commit=`git rev-parse HEAD` -X main.built=`date -u +%s`" \
    -v ./cmd/dxserver

# final stage
FROM alpine:3.24
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /code/dxserver .

LABEL Name=dotagiftx Version=0.25.2
ENTRYPOINT ["./dxserver"]
EXPOSE 80
