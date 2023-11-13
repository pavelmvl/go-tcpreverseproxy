
PROJECT := trp

GOBUILD := go
GOBUILD += build
GOBUILD += -trimpath
GOBUILD += -ldflags='-extldflags=-static -w -s'

GOOS    ?= linux
GOARCH  ?= amd64

OUTPUT  ?= ${PROJECT}.${GOOS}-${GOARCH}

${OUTPUT}: .
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} ${GOBUILD} -o ${OUTPUT} cmd/trp/main.go

all:
	${MAKE} GOOS=linux GOARCH=amd64
	${MAKE} GOOS=linux GOARCH=arm
	${MAKE} GOOS=linux GOARCH=arm64
	${MAKE} GOOS=windows GOARCH=amd64
	${MAKE} GOOS=windows GOARCH=arm64
