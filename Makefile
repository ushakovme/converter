all: build

build:
	CGO_ENABLED=0 go build -o bin/converter cmd/main.go

