.PHONY: build test vet install cross clean

build:
	go build -o bin/keyward ./cmd/keyward

test:
	go test ./... -race

vet:
	go vet ./...

install:
	go install ./cmd/keyward

cross:
	GOOS=darwin  go build ./...
	GOOS=windows go build ./...
	GOOS=linux   go build ./...

clean:
	rm -rf bin dist
