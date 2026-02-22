# build stage
FROM golang:1.26-alpine AS builder
WORKDIR /code

RUN apk add --no-cache git make curl

# download and cache go dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# then copy source code as the last step
COPY . .

RUN make build

# final stage
FROM alpine:3.23
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /code/dxserver .

LABEL Name=dotagiftx Version=0.23.1
ENTRYPOINT exec ./dxserver
EXPOSE 80
