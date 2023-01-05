FROM golang:1.19 AS development
WORKDIR /go/src/github.com/escalopa/gopray/
COPY ./telegram ./telegram
COPY ./pkg ./pkg
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./telegram/cmd/main.go

FROM golang:alpine AS builder
WORKDIR /go/src/github.com/escalopa/gopray/
COPY ./telegram ./telegram
COPY ./pkg ./pkg
COPY go.mod go.sum ./
RUN go build -o /go/bin/gopray ./telegram/cmd

FROM alpine:latest AS production
RUN apk add --no-cache tzdata
COPY --from=builder /go/bin/gopray /go/bin/gopray
ENTRYPOINT ["/go/bin/gopray"]
