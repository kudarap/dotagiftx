# build stage
FROM golang:1.26-alpine AS builder
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
    -ldflags="-X main.tag=`cat VERSION` -X main.commit=`git rev-parse HEAD` -X main.built=`date -u +%s`" \
    -v ./cmd/dxserver

# final stage
FROM alpine:3.23
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /code/dxserver .

LABEL Name=dotagiftx Version=0.23.1
ENTRYPOINT ["./dxserver"]
EXPOSE 80
