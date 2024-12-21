export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: main

fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

gci:
	gci write -s standard -s default -s "prefix(github.com/LoCyan-Team/LoCyanFrpController/)" ./

vet:
	go vet ./...

main:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/controller ./

test: gotest

gotest:
	go test -v --cover ./net/...
	go test -v --cover ./pkg/...
	
clean:
	rm -f ./bin/controller