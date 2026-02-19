.PHONY: build install test clean

build:
	go build -o kipf .

install:
	go install .

test:
	go test -v ./...

clean:
	rm -f kipf
