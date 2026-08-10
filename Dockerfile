# build stage
FROM golang:1.26-alpine AS builder
ARG VERSION
RUN apk add --no-cache git make curl

WORKDIR /code

# download and cache go dependencies
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# then copy source code as the last step
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags="-X main.tag=${VERSION:-`cat VERSION`} -X main.commit=`git rev-parse HEAD` -X main.built=`date -u +%s`" \
    -v ./cmd/dxserver

# final stage
FROM alpine:3.24
RUN apk --no-cache add ca-certificates tzdata

ARG VERSION
COPY --from=builder /code/dxserver .

LABEL Name=dotagiftx Version=${VERSION:-unknown}
ENTRYPOINT ["./dxserver"]
EXPOSE 80
